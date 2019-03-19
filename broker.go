package cuber

import (
	_ "fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

type config struct {
	Pool *redis.Pool
	// Fetch func(queue string) Fetcher
}

var Config *config

func Configure(options map[string]string) {
	if options["server"] == "" {
		panic("Configure requires a 'server' option, which identifies a Redis instance")
	}

	if options["pool"] == "" {
		panic("Configure requires a 'pool' option, which identifies number of connections to keep open with redis")
	}

	var poolSize int
	var err error
	if poolSize, err = strconv.Atoi(options["pool"]); err != nil {
		panic("Option 'pool' must be Integer")
	}
	// fmt.Println(poolSize)

	Config = &config{
		&redis.Pool{
			MaxIdle:     poolSize,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", options["server"])
				if err != nil {
					return nil, err
				}
				if options["password"] != "" {
					if _, err := c.Do("AUTH", options["password"]); err != nil {
						c.Close()
						return nil, err
					}
				}
				if options["database"] != "" {
					if _, err := c.Do("SELECT", options["database"]); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
		// func(queue string) Fetcher {
		// 	return NewFetch(queue, make(chan *Msg), make(chan bool))
		// },
	}
}
