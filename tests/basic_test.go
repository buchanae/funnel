package tests

import (
  "testing"
)

func TestHelloWorld(t *testing.T) {
  id := fun.Run("echo", "hello world")
  fun.WaitForTask(id)
  task := fun.GetTask(id)

  if task.Logs[0].Logs[0].Stdout != "hello world\n" {
    t.Fatal("Missing stdout")
  }
}

// Test that the streaming logs pick up a single character.
// This ensures that the streaming works even when a small
// amount of logs are written (which was once a bug).
func TestSingleCharLog(t *testing.T) {
  id := fun.Run("sh -c 'echo a; fun-wait step-1'")
  fun.WaitForSignal("step-1")
  time.Sleep(time.Millisecond * 100)
  task := fun.GetTask(id)
  if task.Logs[0].Logs[0].Stdout != "a\n" {
    t.Fatal("Missing logs")
  }
}

// Test that port mappings are being logged.
func TestPortLog(t *testing.T) {
  id := fun.Run("fun-wait step-1")
  fun.WaitForSignal("step-1")
  time.Sleep(time.Millisecond * 100)
  task := fun.GetTask(id)
  if task.Logs[0].Logs[0].Ports[0].Host != 5000 {
    t.Fatal("Unexpected port logs")
  }
}

// Test that a completed task cannot change state.
func TestCompleteStateImmutable(t *testing.T) {
  id := fun.Run("echo hello")
  fun.WaitForTask(id)
  fun.CancelTask(id)
  task := fun.GetTask(id)
  if task.State != tes.State_COMPLETE {
    t.Fatal("Unexpected state")
  }
}

// Test canceling a task
func TestCancel(t *testing.T) {
  id := fun.Run("fun-wait step-1", "fun-wait step-2")
  fun.WaitForSignal("step-1")
  fun.WaitForContainer(id + "-0")
  fun.CancelTask(id)
  fun.WaitForContainerStop(task_id + "-0")
  fun.AssertNoContainer(taskID + "-0")
  fun.AssertNoContainer(taskID + "-1")
  task := fun.GetTask(id)
  if task.State != tes.State_CANCELED {
    t.Fatal("Unexpected state")
  }
}

// The task executor logs list should only include entries for steps that
// have been started or completed, i.e. steps that have yet to be started
// won't show up in Task.Logs[0].Logs
func TestExecutorLogLength(t *testing.T) {
  id := fun.Run("fun-wait step-1", "echo done")
  fun.WaitForSignal("step-1")
  task := fun.GetTask(id)
  if len(task.Logs[0].Logs) != 1 {
    t.Fatal("Unexpected executor log count")
  }
}


// There was a bug + fix where the task was being marked complete after
// the first step completed, but the correct behavior is to mark the
// task complete after *all* steps have completed.
func TestMarkCompleteBug(t *testing.T) {
  id := fun.Run("echo step 1", "fun-wait step-2", "echo step 3")
  fun.WaitForSignal("step-2")
  task := fun.GetTask(id)
  if task.State != tes.State_RUNNING {
    t.Fatal("Unexpected task state")
  }
}
