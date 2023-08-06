package main

import (
	"duoduo_fun/api/router"
	"github.com/gin-gonic/gin"
)

func main() {
	//dal.Init()
	r := gin.Default()
	router.Register(r)
	r.Run(":11451")
}
