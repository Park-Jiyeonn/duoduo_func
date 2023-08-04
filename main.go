package main

import (
	"duoduo_fun/api/router"
	"duoduo_fun/dal"
	"github.com/gin-gonic/gin"
)

func main() {
	dal.Init()
	r := new(gin.Engine)
	router.Register(r)
	r.Run(":11451")
}
