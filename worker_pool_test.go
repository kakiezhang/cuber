package cuber

import (
	"testing"
)

type Ctx struct {
	t *testing.T
}

var Ct *testing.T

func myJob(message *Msg) {
	Ct.Logf("myJob: %v", message)
}

func TestWorkerPool(t *testing.T) {
	Ct = t

	Configure(map[string]string{
		"server":   "localhost:6380",
		"database": "5",
		"pool":     "10",
	})

	Process("myQueue", myJob, 10)

	Run()
}
