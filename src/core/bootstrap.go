package core

import (
	"flag"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	VERSION = "1.0.0"
)

type stringsFlag []string

func (o *stringsFlag) String() string {
	return strings.Join([]string(*o), ",")
}

func (o *stringsFlag) Set(value string) error {
	*o = append(*o, value)
	return nil
}

func (o *stringsFlag) Get() []string {
	return []string(*o)
}

func Run() {
	var (
		configPaths stringsFlag
		showHelp    bool
	)
	flag.Var(&configPaths, "c", "configuration yaml.")
	flag.Parse()

	if showHelp {
		log.Infof("version %s\n", VERSION)
		flag.Usage()
		return
	}

	orders := []Order{}

	for _, p := range configPaths {

		var (
			c     *Config
			cData Config
			err   error
		)

		if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
			c, err = OpenOnlineConfig(p)
		} else if filepath.IsAbs(p) {
			c, err = OpenOfflineConfig(p)
		} else {
			log.Errorf("config path (%s) error", p)
			continue
		}

		if err != nil {
			log.Error(err)
			continue
		}
		cData, err = c.Resolve()
		if err != nil {
			log.Error(err)
			return
		}
		orders = append(orders, cData.Orders...)
	}

	if len(orders) == 0 {
		log.Error("no order")
		return
	}

	NewMaster(orders).Run()
}
