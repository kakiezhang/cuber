package cuber

import (
	"sync/atomic"
	"time"
)

type worker struct {
	manager    *manager
	stopChan   chan struct{}
	exitChan   chan struct{}
	currentMsg *Msg
	startedAt  int64
}

func (w *worker) start() {
	go w.work(w.manager.fetch.Messages())
}

func (w *worker) quit() {
	w.stopChan <- struct{}{}
	<-w.exitChan
}

func (w *worker) work(messages chan *Msg) {
	for {
		select {
		case message := <-messages:
			atomic.StoreInt64(&w.startedAt, time.Now().UTC().Unix())
			w.currentMsg = message

			w.process(message)

			atomic.StoreInt64(&w.startedAt, 0)
			w.currentMsg = nil

		case <-w.stopChan:
			w.exitChan <- struct{}{}
			return
		}
	}
}

func (w *worker) process(message *Msg) {
	defer func() {
		recover()
	}()

	w.manager.job(message)
}

func (w *worker) processing() bool {
	return atomic.LoadInt64(&w.startedAt) > 0
}

func newWorker(m *manager) *worker {
	return &worker{m, make(chan struct{}), make(chan struct{}), nil, 0}
}
