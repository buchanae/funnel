package store

import (
  "bytes"
  "context"
  "fmt"
	"github.com/ohsu-comp-bio/funnel/config"
	"github.com/ohsu-comp-bio/funnel/storage"
	"github.com/ohsu-comp-bio/funnel/worker"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
  "github.com/golang/protobuf/jsonpb"
	"github.com/spf13/cobra"
  "io/ioutil"
)

var Cmd = &cobra.Command{
	Use: "store",
}

var configFile string

func init() {
	Cmd.AddCommand(putCmd)
	Cmd.AddCommand(uploadCmd)
  putCmd.Flags().StringVar(&configFile, "config", configFile, "Config file.")
  uploadCmd.Flags().StringVar(&configFile, "config", configFile, "Config file.")
}

var uploadCmd = &cobra.Command{
	Use: "upload",
	RunE: func(cmd *cobra.Command, args []string) error {
    conf := config.DefaultConfig()
    config.ParseFile(configFile, &conf)

    var err error
    store := storage.Storage{}
    store, err = store.WithConfig(conf.Worker.Storage)

    if err != nil {
      return err
    }

    baseDir := args[0]
    taskPath := args[1]
    task := &tes.Task{}
    b, err := ioutil.ReadFile(taskPath)
    if err != nil {
      return err
    }

    r := bytes.NewReader(b)
    err = jsonpb.Unmarshal(r, task)
    if err != nil {
      return err
    }


    m := worker.NewFileMapper(baseDir)
    err = m.MapTask(task)
    if err != nil {
      return err
    }
    fmt.Println(m.Outputs)

    ctx := context.Background()
    out, err := worker.Upload(ctx, m, store)
    fmt.Println(out, err)
    return err
  },
}

var putCmd = &cobra.Command{
	Use: "put",
	RunE: func(cmd *cobra.Command, args []string) error {
    conf := config.DefaultConfig()
    config.ParseFile(configFile, &conf)

    var err error
    store := storage.Storage{}
    store, err = store.WithConfig(conf.Worker.Storage)

    if err != nil {
      return err
    }

    ctx := context.Background()
    ty := tes.FileType_FILE
    if len(args) == 3 && args[2] == "dir" {
      ty = tes.FileType_DIRECTORY
    }
    out, err := store.Put(ctx, args[0], args[1], ty)
    fmt.Println(out, err)
		return err
	},
}
