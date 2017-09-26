package gce

import (
  "encoding/json"
  "fmt"
	"github.com/ohsu-comp-bio/funnel/compute/gce"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/util"
	"github.com/ohsu-comp-bio/funnel/gcp"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/worker"
  "github.com/golang/protobuf/proto"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/spf13/cobra"
  "cloud.google.com/go/pubsub"
  "golang.org/x/net/context"
  "os"
  "time"
  "net/http"
  "sync"
  "syscall"
)

var projectID = "isb-cgc-04-0029"
var subscriptionName = "funnel-workers"

var subnodeCmd = &cobra.Command{
	Use: "subnode",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf := config.DefaultConfig()

		// Check that this is a GCE VM environment.
		// If not, fail.
		meta, merr := gce.LoadMetadata()
		if merr != nil {
			log.Error("Error getting GCE metadata", merr)
    } else {
      log.Info("Loaded GCE metadata")
      log.Debug("GCE metadata", meta)

      var err error
      conf, err = gce.WithMetadataConfig(conf, meta)
      if err != nil {
        return err
      }
		}

		logger.Configure(conf.Scheduler.Node.Logger)
    ctx := context.Background()
    ctx = util.SignalContext(ctx, syscall.SIGINT, syscall.SIGTERM)
    return run(ctx, conf)
	},
}

func run(ctx context.Context, conf config.Config) error {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "ok")
  })

  client, err := pubsub.NewClient(ctx, projectID)
  if err != nil {
    return err
  }

  sub := client.Subscription(subscriptionName)
  sub.ReceiveSettings.MaxOutstandingMessages = 1

  go http.ListenAndServe(":7890", nil)

  return sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    task := tes.Task{}
    err := proto.Unmarshal(m.Data, &task)
    if err != nil {
      m.Nack()
      return
    }

    r, err := worker.NewRPCTaskReader(conf.Worker, task.Id)
    fmt.Println("START", task.Name, task.Id)

	  rtask, terr := r.Task()
    if terr != nil {
      m.Nack()
      return
    }

    if rtask.State == tes.Canceled {
      fmt.Println("SKIPPING CANCELED")
      m.Ack()
      return
    }

    var writer events.Writer

    dswriter, err := gcp.NewDatastoreTaskEventWriter("isb-cgc-04-0029", rtask)
    if err != nil {
      return
    }
    defer dswriter.Flush()
    writer = dswriter

    sdwriter, err := gcp.NewStackdriverEventWriter(ctx, "funnel_tasks", "isb-cgc-04-0029")
    if err == nil {
      writer = events.MultiWriter(dswriter, sdwriter)
    }

    timer := time.NewTimer(time.Second * 5)
    go func() {
      for {
        select {
        case <-timer.C:
          dswriter.Flush()
        case <-ctx.Done():
          return
        }
      }
    }()

    w := worker.DefaultWorker{
      Conf: conf.Worker,
      Mapper: worker.NewFileMapper("/"),
      TaskReader: r,
      Event: writer,
    }
    w.Run(ctx)
    fmt.Println("DONE")

    m.Ack()
   })
}

type emptyWriter struct {}
func (emptyWriter) Write(*events.Event) error {
  return nil
}

type taskReader struct {
  task *tes.Task
  state tes.State
  mtx sync.Mutex
}
func (t *taskReader) id() string {
  t.mtx.Lock()
  defer t.mtx.Unlock()
  if t.task == nil {
    return ""
  }
  return t.task.Id
}
func (t *taskReader) setTask(task *tes.Task) {
  t.mtx.Lock()
  defer t.mtx.Unlock()
  t.task = task
}
func (t *taskReader) cancel() {
  t.mtx.Lock()
  defer t.mtx.Unlock()
  t.state = tes.Canceled
}
func (t *taskReader) Task() (*tes.Task, error) {
  t.mtx.Lock()
  defer t.mtx.Unlock()
  return t.task, nil
}
func (t *taskReader) State() (tes.State, error) {
  t.mtx.Lock()
  defer t.mtx.Unlock()
  return t.state, nil
}

var drainCmd = &cobra.Command{
	Use: "drain",
	RunE: func(cmd *cobra.Command, args []string) error {
    return drain()
	},
}

func drain() error {
  ctx := context.Background()

  client, err := pubsub.NewClient(ctx, projectID)
  if err != nil {
    return err
  }

  sub := client.Subscription(subscriptionName)
  if err != nil {
    return err
  }

  return sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
     fmt.Println("drained")
     m.Ack()
   })
}

var listCmd = &cobra.Command{
	Use: "list",
	RunE: func(cmd *cobra.Command, args []string) error {
    return list()
	},
}

func list() error {
  ctx := context.Background()

  client, err := pubsub.NewClient(ctx, projectID)
  if err != nil {
    return err
  }

  sub := client.Subscription(subscriptionName)
  if err != nil {
    return err
  }

  return sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
     defer m.Nack()

    task := tes.Task{}
    err := proto.Unmarshal(m.Data, &task)
    if err != nil {
      m.Nack()
      return
    }
     log.Info("list", task)
   })
}

var readEventsCmd = &cobra.Command{
  Use: "read-events",
  RunE: func(cmd *cobra.Command, args []string) error {

    dec := json.NewDecoder(os.Stdin)
    l := events.NewLogger("foo")

    for {
      ev := events.Event{}
      if err := dec.Decode(&ev); err != nil {
        return err
      }
      l.Write(&ev)
    }

    return nil
  },
}
