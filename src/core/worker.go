package core

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/tomhjx/porter/model"
)

type worker struct {
	task Task
}

func NewWorker(task Task) *worker {
	return &worker{task: task}
}

func (o *worker) Run() (Task, error) {
	c := make(chan model.Content, 1000)
	m := model.New(o.task.Order.Model)
	m.Read(c)
	os.MkdirAll(o.task.Order.Out, 0777)

	for v := range c {
		log.Infof("model content: %v", v)
		err := saveContent(v, o.task.Order.Out)
		if err != nil {
			log.Error(err)
		}
	}
	o.task.State = TaskSucceedState
	return o.task, nil
}

func saveContent(c model.Content, dest string) error {
	s := fmt.Sprintf("%s|%s", c.Parent, c.Source)
	c.ID = fmt.Sprintf("%x", md5.Sum([]byte(s)))
	fpath := dest + "/" + c.ID + ".json"
	file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	row, _ := json.Marshal(c)
	w.WriteString(string(row))
	w.Flush()
	return nil
}
