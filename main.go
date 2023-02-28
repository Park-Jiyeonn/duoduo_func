package main

import (
	"log"
	base "simple_tiktok/kitex_gen/base/baseservice"
)

func main() {
	svr := base.NewServer(new(BaseServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
