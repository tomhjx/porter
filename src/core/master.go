package core

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

type Order struct {
	Model string `yaml:"model"`
	Out   string `yaml:"out"`
}

type master struct {
	orders []Order
}

func NewMaster(orders []Order) *master {
	return &master{orders: orders}
}

func (o *master) Run() {
	o.dispatch()
	wg.Wait()
}

func (o *master) dispatch() {
	wg.Add(len(o.orders))
	for k, v := range o.orders {
		task := Task{Order: v}
		task.ID = int64(k + 1)
		go func(t Task) {
			defer wg.Done()
			rt, err := NewWorker(t).Run()
			if err != nil {
				log.Error(err.Error())
			}
			log.Infof("work res: %v", rt)
		}(task)
	}
}
