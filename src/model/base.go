package model

import "reflect"

type Content struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Parent    string `json:"parent"`
	Source    string `json:"source"`
	CrawledAt int64  `json:"crawledAt"`
	Sort      int64  `json:"sort"`
}

type Base interface {
	Read(c chan Content)
}

type Factory struct{}

func (o Factory) NewRedisReleaseNotes() *RedisReleaseNotes {
	return &RedisReleaseNotes{}
}

func New(name string) Base {
	f := reflect.ValueOf(Factory{}).MethodByName("New" + name)
	return f.Call([]reflect.Value{})[0].Interface().(Base)
}
