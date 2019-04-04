package cuber

import (
	"sync"
)

type manager struct {
	queue       string
	fetch       Fetcher
	job         jobFunc
	concurrency int
	workers     []*worker
	workersM    *sync.Mutex
	stopChan    chan struct{}
	exitChan    chan struct{}
	*sync.WaitGroup
}

func (m *manager) start() {
	m.Add(1)
	m.loadWorkers()
	go m.manage()
}

func (m *manager) quit() {
	Logger.Println("quitting queue", m.queue, "(waiting for", m.processing(), "/", len(m.workers), "workers).")

	if !m.fetch.Closed() {
		m.fetch.Close()
	}

	m.workersM.Lock()
	for _, worker := range m.workers {
		worker.quit()
	}
	m.workersM.Unlock()

	m.stopChan <- struct{}{}
	<-m.exitChan

	m.reset()

	m.Done()
}

func (m *manager) manage() {
	Logger.Println("processing queue", m.queue, "with", m.concurrency, "workers.")

	go m.fetch.Fetch()

	for {
		select {
		case <-m.stopChan:
			m.exitChan <- struct{}{}
			break
		}
	}
}

func (m *manager) loadWorkers() {
	m.workersM.Lock()
	for i := 0; i < m.concurrency; i++ {
		m.workers[i] = newWorker(m)
		m.workers[i].start()
	}
	m.workersM.Unlock()
}

func (m *manager) processing() (count int) {
	m.workersM.Lock()
	for _, worker := range m.workers {
		if worker.processing() {
			count++
		}
	}
	m.workersM.Unlock()
	return
}

func (m *manager) reset() {
	m.fetch = Config.Fetch(m.queue)
}

func newManager(queue string, job jobFunc, concurrency int) *manager {
	m := &manager{
		queue,
		nil,
		job,
		concurrency,
		make([]*worker, concurrency),
		&sync.Mutex{},
		make(chan struct{}),
		make(chan struct{}),
		&sync.WaitGroup{},
	}

	m.fetch = Config.Fetch(m.queue)

	return m
}
