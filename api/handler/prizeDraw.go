package handler

import (
	"context"
	"duoduo_fun/api/model"
	"fmt"
	"github.com/Park-Jiyeonn/coreRPC/xclient"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var d *xclient.GeeRegistryDiscovery
var xc *xclient.XClient

func PrizeDraw(c *gin.Context) {
	//_, err := strconv.Atoi(c.Query("id"))
	//if err != nil {
	//	c.String(http.StatusBadRequest, "Invalid ID")
	//	return
	//}
	//
	n := 50
	resp := new(model.TreasureResponse)
	//defer func() { _ = xc.Close() }()

	if d == nil || xc == nil {
		d = xclient.NewGeeRegistryDiscovery("http://110.42.239.202:9999/_geerpc_/registry", 0)
		xc = xclient.NewXClient(d, xclient.RandomSelect, nil)
		fmt.Println("新生成了")
	}

	// send request & receive response
	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
	err := xc.Call(ctx, "Treasure.GetTreasure", n, resp)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
