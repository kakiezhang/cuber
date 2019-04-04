package cuber

import (
	"github.com/bitly/go-simplejson"
	_ "reflect"
)

type data struct {
	*simplejson.Json
}

type Msg struct {
	*data
}

func (d *data) ToJson() string {
	json, err := d.Encode()

	if err != nil {
		Logger.Println("ERR: Couldn't generate json from", d, ":", err)
	}

	return string(json)
}

func NewMsg(content string) (*Msg, error) {
	if d, err := newData(content); err != nil {
		return nil, err
	} else {
		return &Msg{d}, nil
	}
}

func newData(content string) (*data, error) {
	if json, err := simplejson.NewJson([]byte(content)); err != nil {
		return nil, err
	} else {
		return &data{json}, nil
	}
}
