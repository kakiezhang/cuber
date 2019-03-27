package cuber

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var Logger WorkerLogger = log.New(os.Stdout, "workers: ", log.Ldate|log.Lmicroseconds)
var managers = make(map[string]*manager)
var access sync.Mutex
var started bool

func Process(queue string, job jobFunc, concurrency int) {
	access.Lock()
	defer access.Unlock()

	managers[queue] = newManager(queue, job, concurrency)
}

func Run() {
	Start()
	go handleSignals()
	waitForExit()
}

func handleSignals() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)

	for sig := range signals {
		switch sig {
		case syscall.SIGINT, syscall.SIGUSR1, syscall.SIGTERM:
			Quit()
		}
	}
}

func ResetManagers() error {
	access.Lock()
	defer access.Unlock()

	if started {
		return errors.New("Cannot reset worker managers while workers are running")
	}

	managers = make(map[string]*manager)

	return nil
}

func Start() {
	access.Lock()
	defer access.Unlock()

	if started {
		return
	}

	startManagers()

	started = true
}

func startManagers() {
	for _, manager := range managers {
		manager.start()
	}
}

func Quit() {
	access.Lock()
	defer access.Unlock()

	if !started {
		return
	}

	quitManagers()
	waitForExit()

	started = false
}

func quitManagers() {
	for _, m := range managers {
		go (func(m *manager) {
			m.quit()
		})(m)
	}
}

func waitForExit() {
	for _, manager := range managers {
		manager.Wait()
	}
}
