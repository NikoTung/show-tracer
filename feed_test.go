package main

import (
    "fmt"
    "github.com/mmcdole/gofeed"
    "testing"
)

func TestFeed(t *testing.T) {
    fp := gofeed.NewParser()
    feed, _ := fp.ParseURL("https://showrss.info/user/261317.rss?magnets=true&namespaces=true&name=null&quality=null&re=null")
    fmt.Println(feed.Title)
    fmt.Println(feed.String())
}

