package worker

import (
	"context"
	"github.com/ohsu-comp-bio/funnel/config"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"io"
)

// TODO document behavior of slow consumer of task log updates

func NewRPCTask(conf config.Config, taskID string) (*RPCTask, error) {
  return &RPCTask{client, taskID}, nil
}

type RPCTask struct {
  client schedClient
  taskID string
}

func (r *RPCTask) Close() {}

func (r *RPCTask) Task() (*tes.Task, error) {
  return r.client.GetTask(context.TODO(), &tes.GetTaskRequest{
    Id: r.taskID,
    View: tes.TaskView_FULL,
  })
}

func (r *RPCTask) State() (*tes.State, error) {
  task, err := r.client.GetTask(context.TODO(), &tes.GetTaskRequest{
    Id: r.taskID,
  })
  if err != nil {
    return nil, err
  }
  return task.State, nil
}

func (r *RPCTask) StartTime(t string) {
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    TaskLog: &tes.TaskLog{
      StartTime: t,
    },
  })
}

func (r *RPCTask) EndTime(t string) {
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    TaskLog: &tes.TaskLog{
      EndTime: t,
    },
  })
}

func (r *RPCTask) Outputs(f []string) {
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    TaskLog: &tes.TaskLog{
      Outputs: f,
    },
  })
}

func (r *RPCTask) Metadata(m map[string]string) {
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    TaskLog: &tes.TaskLog{
      Metadata: m,
    },
  })
}

func (r *RPCTask) Running() {
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    State: tes.State_RUNNING,
  })
}

func (r *RPCTask) Result(err error) {
  var state tes.State
  if err == nil {
    state = tes.State_COMPLETE
  } else {
    state = tes.State_ERROR
  }
  // TODO SYSTEM_ERROR

  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    State: state,
  })
}

func (r *RPCTask) ExecutorStartTime(i int, t string) {
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    ExecutorIndex: i,
    ExecutorLog: &tes.ExecutorLog{,
      StartTime: t,
    },
  })
}

func (r *RPCTask) ExecutorEndTime(i int, t string) {
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    ExecutorIndex: i,
    ExecutorLog: &tes.ExecutorLog{,
      EndTime: t,
    },
  })
}

func (r *RPCTask) ExecutorExitCode(i int, x int) {
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    ExecutorIndex: i,
    ExecutorLog: &tes.ExecutorLog{,
      ExitCode: int32(x),
    },
  })
}

func (r *RPCTask) ExecutorPorts(i int, ports []*tes.Ports) {
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    ExecutorIndex: i,
    ExecutorLog: &tes.ExecutorLog{,
      Ports: ports,
    },
  })
}

func (r *RPCTask) ExecutorHostIP(i int, ip string) {
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    ExecutorIndex: i,
    ExecutorLog: &tes.ExecutorLog{,
      HostIp: ip,
    },
  })
}

func (r *RPCTask) ExecutorStdout(i int) io.Writer {
  // tailer
    // TODO
    //Stdout: io.MultiWriter(stdout, log.Stdout())
    //Stderr: io.MultiWriter(stderr, log.Stderr())
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    ExecutorIndex: i,
    ExecutorLog: &tes.ExecutorLog{,
      Stdout: "",
    },
  })
}

func (r *RPCTask) ExecutorStderr(i int) io.Writer {
  // tailer
  r.client.UpdateTask(context.TODO(), &pbf.UpdateTaskRequest{
    Id: r.taskID,
    ExecutorIndex: i,
    ExecutorLog: &tes.ExecutorLog{
      Stderr: "",
    },
  })
}
