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

func (t *tes.Task) GetTaskLog(i int) *tes.TaskLog {

	// Grow slice length if necessary
	if len(t.Logs) <= i {
		desired := i + 1
		t.Logs = append(t.Logs, make([]*TaskLog, desired-len(t.Logs))...)
	}

	if t.Logs[i] == nil {
		t.Logs[i] = &TaskLog{}
	}

	return t.Logs[i]
}

func (t *Task) GetExecLog(attempt int, i int) *ExecutorLog {
	tl := getTaskLog(t, 0)

	// Grow slice length if necessary
	if len(tl.Logs) <= i {
		desired := i + 1
		tl.Logs = append(tl.Logs, make([]*ExecutorLog, desired-len(tl.Logs))...)
	}

	if tl.Logs[i] == nil {
		tl.Logs[i] = &ExecutorLog{}
	}

	return tl.Logs[i]
}
