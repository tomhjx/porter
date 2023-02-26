package model

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

type RedisReleaseNotes struct {
}

func (o *RedisReleaseNotes) Read(contentNotify chan Content) {

	defer close(contentNotify)

	source := "https://redis.io/download/"

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if !strings.Contains(strings.ToLower(e.Text), "release notes") {
			return
		}
		url := e.Attr("href")
		if !strings.Contains(url, "//raw.githubusercontent.com/redis/redis/") {
			return
		}
		ver := strings.Split(url, "/")[5]

		log.Infof("%s, %s, %s", e.Text, url, ver)

		ctx := Content{}
		ctx.Title = fmt.Sprintf("Redis %s Release Notes", ver)
		ctx.Source = url
		ctx.Parent = source
		ctx.CrawledAt = time.Now().Unix()
		resp, err := http.Get(url)
		if err != nil {
			log.Errorf("%s, err: %s", url, err)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("%s, body read err: %s", url, err)
			return
		}
		ctx.Body = string(body)
		contentNotify <- ctx
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(source)
	c.Wait()
}
