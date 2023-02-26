package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tomhjx/porter/core"
	"gopkg.in/yaml.v3"
)

const (
	VERSION = "1.0.0"
)

type stringsFlag []string

func (i *stringsFlag) String() string {
	return strings.Join([]string(*i), ",")
}

func (i *stringsFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *stringsFlag) Get() []string {
	return []string(*i)
}

func main() {
	var (
		srcps    stringsFlag
		showHelp bool
	)
	flag.Var(&srcps, "s", "source's clashx configuration yaml url.")
	flag.Parse()

	allOrders := []core.Order{}

	if showHelp {
		fmt.Printf("version %s\n", VERSION)
		flag.Usage()
		return
	}
	for _, srcp := range srcps {
		c, err := ioutil.ReadFile(srcp)
		if err != nil {
			log.Error(err)
			continue
		}
		orders := []core.Order{}
		err = yaml.Unmarshal(c, &orders)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("%v", orders)
		allOrders = append(allOrders, orders...)
	}
	core.NewMaster(allOrders).Run()
}
