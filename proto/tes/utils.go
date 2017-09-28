package tes

// RunnableState returns true if the state is RUNNING or INITIALIZING
func RunnableState(s State) bool {
	return s == State_INITIALIZING || s == State_RUNNING
}

// TerminalState returns true if the state is COMPLETE, ERROR, SYSTEM_ERROR, or CANCELED
func TerminalState(s State) bool {
	return s == State_COMPLETE || s == State_ERROR || s == State_SYSTEM_ERROR ||
		s == State_CANCELED
}

func (t *Task) GetTaskLog(i int) *TaskLog {

	// Grow slice length if necessary
  for j := len(t.Logs); j <= i; j++ {
    t.Logs = append(t.Logs, &TaskLog{})
  }

	return t.Logs[i]
}

func (t *Task) GetExecLog(attempt int, i int) *ExecutorLog {
	tl := t.GetTaskLog(attempt)

	// Grow slice length if necessary
  for j := len(tl.Logs); j <= i; j++ {
    tl.Logs = append(tl.Logs, &ExecutorLog{})
  }

	return tl.Logs[i]
}
