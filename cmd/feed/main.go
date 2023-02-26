package main

import (
	"log"
	feed "simple_tiktok/kitex_gen/feed/feedservice"
)

func main() {
	svr := feed.NewServer(new(FeedServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
