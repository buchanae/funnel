package worker

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/events"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"github.com/ohsu-comp-bio/funnel/server/dynamodb"
	"github.com/ohsu-comp-bio/funnel/server/elastic"
	"github.com/ohsu-comp-bio/funnel/server/mongodb"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/worker"
	"golang.org/x/oauth2/google"
	gsstorage "google.golang.org/api/storage/v1"
	"io/ioutil"
	"strings"
)

type GoogleStorageTaskReader struct {
	svc    *gsstorage.Service
	taskID string
}

func NewGoogleStorageTaskReader(taskID string) (*GoogleStorageTaskReader, error) {
	ctx := context.Background()
	// Pull the information (auth and other config) from the environment,
	// which is useful when this code is running in a Google Compute instance.
	client, err := google.DefaultClient(ctx, gsstorage.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	svc, cerr := gsstorage.New(client)
	if cerr != nil {
		return nil, cerr
	}
	return &GoogleStorageTaskReader{svc, taskID}, nil
}
func (g *GoogleStorageTaskReader) Task() (*tes.Task, error) {
	bucket, path := parse(g.taskID)
	call := g.svc.Objects.Get(bucket, path)
	resp, derr := call.Download()
	if derr != nil {
		return nil, fmt.Errorf("can't get task: %s", derr)
	}
	var task tes.Task

	b, cerr := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if cerr != nil {
		return nil, fmt.Errorf("can't get task: %s", derr)
	}

	err := jsonpb.UnmarshalString(string(b), &task)
	if err != nil {
		return nil, fmt.Errorf("can't get task: %s", derr)
	}
	task.Id = tes.GenerateID()
	return &task, nil
}
func (g *GoogleStorageTaskReader) State() (tes.State, error) {
	return tes.Unknown, nil
}
func parse(rawurl string) (string, string) {
	path := strings.TrimPrefix(rawurl, "gs://")
	split := strings.SplitN(path, "/", 2)
	bucket := split[0]
	key := split[1]
	return bucket, key
}

// Run configures and runs a Worker
func Run(ctx context.Context, conf config.Worker, taskID string, log *logger.Logger) error {
	log.Debug("Run Worker", "config", conf, "taskID", taskID)

	var err error
	var db tes.ReadOnlyServer
	var reader worker.TaskReader
	var writer events.Writer

	switch conf.TaskReader {
	case "file":
		reader, err = worker.NewFileTaskReader(taskID)
	case "gs":
		reader, err = NewGoogleStorageTaskReader(taskID)
	case "rpc":
		reader, err = worker.NewRPCTaskReader(conf.TaskReaders.RPC, taskID)
	case "dynamodb":
		db, err = dynamodb.NewDynamoDB(conf.TaskReaders.DynamoDB)
	case "elastic":
		db, err = elastic.NewElastic(ctx, conf.EventWriters.Elastic)
	case "mongodb":
		db, err = mongodb.NewMongoDB(conf.TaskReaders.MongoDB)
	default:
		err = fmt.Errorf("unknown TaskReader")
	}
	if err != nil {
		return fmt.Errorf("failed to instantiate TaskReader: %v", err)
	}

	if reader == nil {
		reader = worker.NewGenericTaskReader(db.GetTask, taskID)
	}

	writers := []events.Writer{}
	for _, w := range conf.ActiveEventWriters {
		switch w {
		case "log":
			writer = &events.Logger{Log: log}
		case "rpc":
			writer, err = events.NewRPCWriter(conf.EventWriters.RPC)
		case "dynamodb":
			writer, err = dynamodb.NewDynamoDB(conf.EventWriters.DynamoDB)
		case "elastic":
			writer, err = elastic.NewElastic(ctx, conf.EventWriters.Elastic)
		case "mongodb":
			writer, err = mongodb.NewMongoDB(conf.EventWriters.MongoDB)
		case "kafka":
			k, kerr := events.NewKafkaWriter(conf.EventWriters.Kafka)
			defer k.Close()
			err = kerr
			writer = k
		default:
			err = fmt.Errorf("unknown EventWriter")
		}
		if err != nil {
			return fmt.Errorf("failed to instantiate EventWriter: %v", err)
		}
		writers = append(writers, writer)
	}

	m := events.MultiWriter(writers)
	ew := &events.ErrLogger{Writer: &m, Log: log}

	w := &worker.DefaultWorker{
		Conf:        conf,
		Store:       storage.Storage{},
		TaskReader:  reader,
		EventWriter: ew,
	}
	// TODO doesn't pass container exit code through funnel/root.go and cobra.
	return w.Run(ctx)
}
