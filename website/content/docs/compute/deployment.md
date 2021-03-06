---
title: Deploying a cluster
menu:
  main:
    parent: Compute
    weight: -50
---

# Deploying a cluster

This guide describes the basics of starting a cluster of Funnel nodes. 
This guide is a work in progress.

A node is a service
which runs on each machine in a cluster. The node connects to the Funnel server and reports
available resources. The Funnel scheduler process assigns tasks to nodes. When a task is
assigned, a node will start a worker process. There is one worker process per task.

Nodes aren't always required. In some cases it makes sense to rely on an existing,
external system for scheduling tasks and managing cluster resources, such as AWS Batch,
HTCondor, Slurm, Grid Engine, etc. Funnel provides integration with
these services without using nodes or the scheduler.

### Usage

Nodes are available via the `funnel node` command. To start a node, run
```
funnel node run --config node.config.yml
```

To activate the Funnel scheduler, use the `manual` backend in the config.

The available scheduler and node config:
```
# Activate the Funnel scheduler.
Backend: manual

Scheduler:
  # How often to run a scheduler iteration.
  # In nanoseconds.
  ScheduleRate: 1000000000 # 1 second

  # How many tasks to schedule in one iteration.
  ScheduleChunk: 10

  # How long to wait between updates before marking a node dead.
  # In nanoseconds.
  NodePingTimeout: 60000000000 # 1 minute

  # How long to wait for a node to start, before marking the node dead.
  # In nanoseconds.
  NodeInitTimeout: 300000000000 # 5 minutes

  # Node config.
  Node:
    # If empty, a node ID will be automatically generated using the hostname.
    ID: ""

    # Files created during processing will be written in this directory.
    WorkDir: ./funnel-work-dir

    # If the node has been idle for longer than the timeout, it will shut down.
    # -1 means there is no timeout. 0 means timeout immediately after the first task.
    Timeout: -1

    # A Node will automatically try to detect what resources are available to it. 
    # Defining Resources in the Node configuration overrides this behavior.
    Resources:
      # CPUs available.
      # Cpus: 0
      # RAM available, in GB.
      # RamGb: 0.0
      # Disk space available, in GB.
      # DiskGb: 0.0

    # For low-level tuning.
    # How often to sync with the Funnel server.
    # In nanoseconds.
    UpdateRate: 5000000000 # 5 seconds

    # RPC timeout for update/sync call.
    # In nanoseconds.
    UpdateTimeout: 1000000000 # 1 second

    Logger:
      # Logging levels: debug, info, error
      Level: info
      # Write logs to this path. If empty, logs are written to stderr.
      OutputFile: ""
```

### Known issues

The config uses nanoseconds for duration values. See [issue #342](https://github.com/ohsu-comp-bio/funnel/issues/342).
