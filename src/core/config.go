package core

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Orders []Order
	data   []byte
}

func OpenOfflineConfig(file string) (c *Config, err error) {
	log.Infof("open local config: %s", file)
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return &Config{data: d}, nil
}

func OpenOnlineConfig(url string) (c *Config, err error) {
	log.Infof("open remote config: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 porter"+VERSION)
	httpc := &http.Client{}
	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("online request returned code: %d", resp.StatusCode)
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Config{data: d}, nil
}

func (o *Config) Resolve() (Config, error) {
	c := Config{}
	err := yaml.Unmarshal(o.data, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}
