package worker

type FileTaskReader struct {}
func (b *FileTaskReader) Task() *tes.Task {
  return b.task
}
func (f *FileTaskReader) State() tes.State {
}
