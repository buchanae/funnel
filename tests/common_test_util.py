from __future__ import print_function

import logging
import os
import py_tes
import shutil
import signal
import subprocess
import tempfile
import time
import unittest
import requests
import polling
import yaml
import docker


S3_ENDPOINT = "localhost:9000"
S3_ACCESS_KEY = "AKIAIOSFODNN7EXAMPLE"
S3_SECRET_KEY = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
BUCKET_NAME = "tes-test"
WORK_DIR = os.path.join(os.path.dirname(os.path.dirname(__file__)),
                        "test_work")


def popen(*args, **kwargs):
    kwargs['preexec_fn'] = os.setsid
    return subprocess.Popen(*args, **kwargs)


def kill(p):
    try:
        os.killpg(os.getpgid(p.pid), signal.SIGTERM)
        p.wait()
    except OSError:
        pass


def get_abspath(path):
    return os.path.join(os.path.dirname(__file__), path)


def which(file):
    for path in os.environ["PATH"].split(":"):
        p = os.path.join(path, file)
        if os.path.exists(p):
            return p


def temp_config(config):
    configFile = tempfile.NamedTemporaryFile(delete=False)
    yaml.dump(config, configFile)
    return configFile


def config_seconds(sec):
    # The funnel config is currently parsed as nanoseconds
    # this helper makes that manageale
    return int(sec * 1000000000)


task_server = None
tempdir = tempfile.mkdtemp(prefix="funnel-py-tests-")
storage_dir = os.path.abspath(tempdir + ".storage")
work_dir = os.path.abspath(tempdir + ".work-dir")


def teardown_package():
    if task_server is not None:
        kill(task_server)

signal.signal(signal.SIGINT, teardown_package)


def setup_package():
    global task_server
    os.mkdir(storage_dir)
    os.mkdir(work_dir)
    # Build server config file (YAML)
    rate = config_seconds(0.05)
    configFile = temp_config({
        "HostName": "localhost",
        "HTTPPort": "8000",
        "RPCPort": "9090",
        "WorkDir": work_dir,
        "Storage": [{
            "Local": {
                "AllowedDirs": [storage_dir]
            },
            "S3": {
                "Endpoint": S3_ENDPOINT,
                "Key": S3_ACCESS_KEY,
                "Secret": S3_SECRET_KEY,
            }
        }],
        "LogLevel": "info",
        "Worker": {
            "Timeout": -1,
            "StatusPollRate": rate,
            "LogUpdateRate": rate,
            "NewJobPollRate": rate,
            "UpdateRate": rate,
            "TrackerRate": rate
        },
        "ScheduleRate": rate,
    })
    # Start server
    cmd = ["funnel", "server", "--config", configFile.name]
    logging.info("Running %s" % (" ".join(cmd)))
    task_server = popen(cmd)
    time.sleep(1)


class SimpleServerTest(unittest.TestCase):

    def setUp(self):
        self.storage_dir = storage_dir
        self.tes = py_tes.TES("http://localhost:8000")

    def storage_path(self, *args):
        return os.path.join(self.storage_dir, *args)

    def copy_to_storage(self, path):
        dst = os.path.join(self.storage_dir, os.path.basename(path))
        shutil.copy(path, dst)
        return os.path.basename(path)

    def get_from_storage(self, loc):
        dst = os.path.join(self.storage_dir, loc)
        return dst

    def wait_for_container(self, name, timeout=5):
        dclient = docker.from_env()

        def on_poll():
            try:
                dclient.containers.get(name)
                return True
            except BaseException:
                return False
        polling.poll(on_poll, timeout=timeout, step=0.1)

    def wait_for_container_stop(self, name, timeout=5):
        dclient = docker.from_env()

        def on_poll():
            try:
                dclient.containers.get(name)
                return False
            except BaseException:
                return True
        polling.poll(on_poll, timeout=timeout, step=0.1)

    def wait(self, key, timeout=5):
        """
        Waits for tes-wait to return <key>
        """
        def on_poll():
            try:
                r = requests.get("http://127.0.0.1:5000/")
                return r.status_code == 200 and r.text == key
            except requests.ConnectionError:
                return False

        polling.poll(on_poll, timeout=timeout, step=0.1)

    def resume(self):
        """
        Continue from tes-wait
        """
        requests.get("http://127.0.0.1:5000/shutdown")


class S3ServerTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        cls.output_dir = os.path.join(tempdir + ".s3-output")
        os.mkdir(cls.output_dir)
        cls.dir_name = os.path.basename(tempdir)

        # start s3 server
        cmd = [
            which("docker"),
            "run", "-p", "9000:9000",
            "--rm",
            "--name", "tes_minio_test",
            "-e", "MINIO_ACCESS_KEY=%s" % (S3_ACCESS_KEY),
            "-e", "MINIO_SECRET_KEY=%s" % (S3_SECRET_KEY),
            "-v", "%s:/export" % (storage_dir),
            "minio/minio", "server", "/export"
        ]

        logging.info("Running %s" % (" ".join(cmd)))
        cls.s3_server = popen(cmd)

        # TES client
        cls.tes = py_tes.TES(
            "http://localhost:8000",
            "http://" + S3_ENDPOINT,
            S3_ACCESS_KEY,
            S3_SECRET_KEY
        )

        if not cls.tes.bucket_exists(BUCKET_NAME):
            cls.tes.make_bucket(BUCKET_NAME)

    @classmethod
    def teardownClass(cls):
        if cls.s3_server is not None:
            kill(cls.s3_server)
            cmd = ["docker", "kill", "tes_minio_test"]
            logging.info("Running %s" % (" ".join(cmd)))

            popen(cmd).communicate()

            cmd = ["docker", "rm", "-fv", "tes_minio_test"]
            logging.info("Running %s" % (" ".join(cmd)))

            popen(cmd).communicate()

    def get_storage_url(self, path):
        dstpath = "s3://%s/%s" % (
            BUCKET_NAME, os.path.join(self.dir_name, os.path.basename(path))
        )
        return dstpath

    def copy_to_storage(self, path):
        dstpath = "s3://%s/%s" % (
            BUCKET_NAME, os.path.join(self.dir_name, os.path.basename(path))
        )
        print("uploading:", dstpath)
        self.tes.upload_file(path, dstpath)
        return dstpath

    def get_from_storage(self, loc):
        dst = os.path.join(self.output_dir, os.path.basename(loc))
        print("Downloading %s to %s" % (loc, dst))
        self.tes.download_file(dst, loc)
        return dst
