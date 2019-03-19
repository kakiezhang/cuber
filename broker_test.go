package cuber

import (
	"testing"
)

var redis_config = map[string]string{
	// location of redis instance
	"server": "localhost:6379",
	// instance of the database
	"database": "0",
	// number of connections to keep open with redis
	"pool": "30",
	// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
	"process": "1",
}

func TestConfig(t *testing.T) {
	Configure(redis_config)
	t.Logf("%v", Config.Pool)
}
