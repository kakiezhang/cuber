package cuber

import (
	_ "reflect"
	_ "strconv"
	"time"

	_ "github.com/bitly/go-simplejson"
	"github.com/gomodule/redigo/redis"
)

type Fetcher interface {
	Queue() string
	Fetch()
	Messages() chan *Msg
	Close()
	Closed() bool
}

type fetch struct {
	queue      string
	messages   chan *Msg
	stopChan   chan struct{}
	exitChan   chan struct{}
	closedChan chan struct{}
}

func NewFetch(queue string, messages chan *Msg) Fetcher {
	return &fetch{
		queue,
		messages,
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}
}

func (f *fetch) Queue() string {
	return f.queue
}

func (f *fetch) Fetch() {
	go func() {
		for {
			// f.Close() has been called
			if f.Closed() {
				break
			}

			f.tryFetchMessage()
		}
	}()

	for {
		select {
		case <-f.stopChan:
			// Stop the redis-polling goroutine
			close(f.closedChan)
			// Signal to Close() that the fetcher has stopped
			close(f.exitChan)
			break
		}
	}
}

func (f *fetch) tryFetchMessage() {
	conn := Config.Pool.Get()
	defer conn.Close()

	message, err := redis.Values(conn.Do("blpop", f.queue, 0))

	if err != nil {
		// If redis returns null, the queue is empty. Just ignore the error.
		if err.Error() != "redigo: nil returned" {
			Logger.Println("ERR: ", err)
			time.Sleep(1 * time.Second)
		}
	} else {
		// Logger.Printf("msg111: %v, %s", message[1], reflect.TypeOf(message[1]))

		v, ok := message[1].([]uint8)
		if ok {
			bs := B2S(v)
			Logger.Printf("bs: %s\n", bs)
			f.sendMessage(bs)
		} else {
			Logger.Printf("bs err\n")
		}
		// json, _ := simplejson.NewJson([]byte(message[1]))
		// Logger.Printf("json111: %v\n", json)
	}
}

func B2S(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

func (f *fetch) sendMessage(message string) {
	msg, err := NewMsg(message)

	if err != nil {
		Logger.Println("ERR: Couldn't create message from", message, ":", err)
		return
	}

	f.Messages() <- msg
}

func (f *fetch) Messages() chan *Msg {
	return f.messages
}

func (f *fetch) Close() {
	f.stopChan <- struct{}{}
	<-f.exitChan
}

func (f *fetch) Closed() bool {
	select {
	case <-f.closedChan:
		return true
	default:
		return false
	}
}
