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
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strings"
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
	_, err = crons.AddFunc("@every 10m", func() {
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

	sort.Slice(feed.Items, func(i, j int) bool {
		itemi := feed.Items[i]
		itemj := feed.Items[j]

		//desc
		return itemi.PublishedParsed.After(*itemj.PublishedParsed)
	})


	for _, item := range feed.Items {
		if updateTime.Before(*item.PublishedParsed) {
			fmt.Println("Download ", item.Title)
			go download(item.GUID, item.Link, config.Api)
			go sendToTelegram(item.Title, config)
		}
	}

	updateTime = *feed.Items[0].PublishedParsed

}

func fetch(rss string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rss)

	return feed, err
}

func download(guid, link, api string) {

	b := NewAria2(guid, link)

	body, err := json.Marshal(b)

	bb := bytes.NewReader(body)

	req, err := http.NewRequest("POST", api, bb)
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

func sendToTelegram(name string, config *Config)  {
	params := url.Values{}
	params.Add("text", fmt.Sprintf("Show %s submit to download.", name))
	params.Add("chat_id", config.TelegramChatId)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.TelegramToken), body)
	if err != nil {
		// handle err
		fmt.Println("create request error")
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Send message to telegram error")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

}
