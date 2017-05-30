package e2e

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Start minio
	dockerPath, _ := exec.LookPath("docker")
	args := []string{dockerPath, "run",
		"-p", "9999:9999",
		"--rm",
		"--name", "fun-minio-test",
		"-e", "MINIO_ACCESS_KEY=" + fun.Conf.Storage[1].S3.Key,
		"-e", "MINIO_SECRET_KEY=" + fun.Conf.Storage[1].S3.Secret,
		"-v", fun.StorageDir + ":/export",
		"minio/minio", "server", "/export",
	}
	log.Debug("Start minio", strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	// Run the tests
	e := m.Run()

	// Clean up minio
	log.Debug("Stop minio")
	exec.Command(dockerPath, "rm", "-fv", "fun-minio-test").Run()

	// Finish
	os.Exit(e)
}

func TestS3(t *testing.T) {
	id := fun.Run(`
    --cmd "sh -c 'echo foo > $out'"
    -o out=s3://bkt/test_output
  `)
	fun.Wait(id)
}

/*
S3_SECRET_KEY = BUCKET_NAME = "tes-test"

func TestS3(t *testing.T) {
  fun.Storage.Put("s3://bkt/test_input", "test_input", tes.FileType_FILE)
  fun.Storage.Get("s3://bkt/test_output", "test_output", tes.FileType_FILE)
  if readFile("test_output") != "hello-s3" {
    t.Fatal("unexpected s3 output content")
  }
}
*/
