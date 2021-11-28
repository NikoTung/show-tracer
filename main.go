package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mmcdole/gofeed"
	cron "github.com/robfig/cron/v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configFile string

var updateTime = time.Now().UTC()

func init() {
	flag.StringVar(&configFile, "f", "", "the config file location.")

}

func main() {
	flag.Parse()
	if len(configFile) == 0 {
		flag.PrintDefaults()
		os.Exit(2)
	}

	c, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("config file error ", err)
		os.Exit(2)
	}
	var config Config
	err = json.Unmarshal(c, &config)
	if err != nil {
		fmt.Println("config format error ", err)
		os.Exit(2)
	}

	_, err = fetch(config.Rss)

	if err != nil {
		fmt.Println("can not parse the rss url ", err)
		os.Exit(2)
	}

	crons := cron.New()
	_, err = crons.AddFunc("@every 30s", func() {
		update(&config)
	})

	if err != nil {
		fmt.Println("add cron failed ", err)
		os.Exit(2)
	}
	crons.Start()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-sigs
		done <- true
	}()

	<-done
	crons.Stop()

}

func update(config *Config) {
	feed, err := fetch(config.Rss)
	if err != nil {
		fmt.Println("update rss error ", err)
		return
	}

	updateTime = time.Now().UTC()

	for _, item := range feed.Items {
		if updateTime.Before(*item.PublishedParsed) {
			go download(item.GUID, item.Link)
		}
	}

}

func fetch(rss string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rss)

	return feed, err
}

func download(guid, link string) {

	b := NewAria2(guid, link)

	body, err := json.Marshal(b)

	bb := bytes.NewReader(body)

	req, err := http.NewRequest("POST", "http://localhost:6800/jsonrpc", bb)
	if err != nil {
		fmt.Println("new request error ", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("submit download error ", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Close response error ", err)
		}
	}(resp.Body)
}
