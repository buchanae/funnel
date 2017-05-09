package run

import (
	"github.com/ohsu-comp-bio/funnel/cmd/client"
	"github.com/ohsu-comp-bio/funnel/proto/tes"
	"sync"
)

type taskGroup struct {
	wg  sync.WaitGroup
	err chan error
}

func (tg *taskGroup) runTask(t *tes.Task, cli *client.Client, wait bool, waitFor []string) {
	if tg.err == nil {
		tg.err = make(chan error)
	}

	tg.wg.Add(1)
	go func() {
		err := runTask(t, cli, wait, waitFor)
		if err != nil {
			tg.err <- err
		}
		tg.wg.Done()
	}()
}

func (tg *taskGroup) wait() error {
	done := make(chan struct{})
	go func() {
		tg.wg.Wait()
		close(done)
	}()

	select {
	case err := <-tg.err:
		return err
	case <-done:
		return nil
	}
}
