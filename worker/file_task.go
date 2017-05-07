package worker

import (
	"github.com/ohsu-comp-bio/funnel/proto/tes"
)

func NewFileTask(conf config.Config, taskID string) (*FileTask, error) {
}

type FileTask struct {
  taskID string
}

func (r *FileTask) Close() {}

func (r *FileTask) Task() (*tes.Task, error) {
}

func (r *FileTask) State() (*tes.State, error) {
}

func (r *FileTask) StartTime(t string) {
}

func (r *FileTask) EndTime(t string) {
}

func (r *FileTask) Outputs(f []string) {
}

func (r *FileTask) Metadata(m map[string]string) {
}

func (r *FileTask) Running() {
}

func (r *FileTask) Result(err error) {
}

func (r *FileTask) ExecutorStartTime(i int, t string) {
}

func (r *FileTask) ExecutorEndTime(i int, t string) {
}

func (r *FileTask) ExecutorExitCode(i int, x int) {
}

func (r *FileTask) ExecutorHostIP(i int, ip string) {
}

func (r *FileTask) ExecutorStdout(i int) io.Writer {
}

func (r *FileTask) ExecutorStderr(i int) io.Writer {
}
