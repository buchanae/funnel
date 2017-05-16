// Test that a file can be passed as an input and output.
func TestFileMount(t *testing.T) {
  fun.WriteFile("test_in", "hello")
  id := fun.Run("sh -c 'cat $in > $out'", fun.Args(`
    -i in=test_in
    -o out=test_out
  `))
  task := fun.WaitForTask(id)
  c := fun.ReadFile("test_out")
  if c != "hello" {
    t.Fatal("Unexpected output file")
  }
}

// Test that the local storage system hard links output files.
func TestLocalFilesystemHardLink(t *testing.T) {
  fun.WriteFile("test_in", "hello")
  id := fun.Run("sh -c 'cat $in > $out'", fun.Args(`
    -i in=test_in
    -o out=test_out
  `))
}

// Test using a symlink as an input file.
func TestSymlinkInput(t *testing.T) {
  fun.WriteFile("test_in", "hello")
  id := fun.Run("sh -c 'cat $in > $out'", fun.Args(`
    -i in=test_in
    -o out=test_out
  `))
  if task.State != tes.State_COMPLETE {
    t.Fatal("Expected success on symlink input")
  }
}

// Test using a broken symlink as an input file.
func TestBrokenSymlinkInput(t *testing.T) {
  fun.WriteFile("test_in", "hello")
  id := fun.Run("sh -c 'cat $in > $out'", fun.Args(`
    -i in=test_in
    -o out=test_out
  `))
  task := fun.WaitForTask(id)
  if task.State != tes.State_ERROR {
    t.Fatal("Expected error on broken symlink input")
  }
}

/*
  Test the case where a container creates a symlink in an output path.
  From the view of the host system where Funnel is running, this creates
  a broken link, because the source of the symlink is a path relative
  to the container filesystem.

  Funnel can fix some of these cases using volume definitions, which
  is being tested here.
*/
func TestSymlinkOutput(t *testing.T) {
  id := fun.Run("sh -c 'echo foo > $dir/foo && ln-s $dir/foo $dir/sym && ln -s $dir/foo $sym'",
  fun.Args(`
    -o sym=out-sym
    -O dir=out-dir
  `))
  task := fun.WaitForTask(id)

  if task.State != tes.State_COMPLETE {
    t.Fatal("expected success on symlink output")
  }

  if fun.ReadFile("out-dir/foo") != "foo\n" {
    t.Fatal("unexpected out-dir/foo content")
  }

  if fun.ReadFile("out-sym") != "foo\n" {
    t.Fatal("unexpected out-sym content")
  }

  if fun.ReadFile("out-dir/sym") != "foo\n" {
    t.Fatal("unexpected out-dir/sym content")
  }
}

func TestS3(t *testing.T) {
  fun.WriteFile("test_input", "hello-s3")
  fun.Storage.Put("s3://bkt/test_input", "test_input", tes.FileType_FILE)
  id := fun.Run(`sh -c "cat $in > $out"`, fun.Args(`
    -i in=s3://bkt/test_input
    -o out=s3://bkt/test_output
  `))
  fun.WaitForTask(id)
  fun.Storage.Get("s3://bkt/test_output", "test_output", tes.FileType_FILE)
  if fun.ReadFile("test_output") != "hello-s3" {
    t.Fatal("unexpected s3 output content")
  }
}
