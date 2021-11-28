package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"testing"
)

func TestFeed(t *testing.T) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("http://rarbg.to/rssdd.php?categories=4;44;51;54;23;40;14;45;52;18;25;32;48;47;42;41;27;33;17;50;46;49;28;53")
	fmt.Println(feed.Title)
	fmt.Println(feed.String())
}
