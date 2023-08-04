package handler

import (
	"duoduo_fun/api/model"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
)

func generate(treasure []model.Treasure) []model.Treasure {
	n := len(treasure)
	for i := 0; i < n; i++ {
		treasure[i].Value = rand.Intn(50) + 1
		treasure[i].Value = rand.Intn(50) + 1
	}
	return treasure
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
func checkValid(treasure []model.Treasure, capacity, limit int) bool {
	n := len(treasure)
	dp := make([]int, capacity+1)
	for i := 0; i < n; i++ {
		for j := capacity; j >= treasure[i].Weight; j-- {
			dp[j] = max(dp[j], dp[j-treasure[i].Weight]+treasure[i].Value)
		}
	}
	return dp[capacity] <= limit
}

func PrizeDraw(c *gin.Context) {
	_, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID")
	}

	n := 50
	treasure := make([]model.Treasure, n)
	capacity := 0
	limit := 0
	generateSuccess := false

	for i := 0; i < 1000; i++ {
		treasure = generate(treasure)
		capacity = rand.Intn(1000) + 1
		limit = rand.Intn(500) + 1
		if checkValid(treasure, capacity, limit) {
			generateSuccess = true
			break
		}
	}

	if !generateSuccess {
		treasure[0].Value, treasure[0].Weight = limit, capacity
		for i := 1; i < 50; i++ {
			treasure[i].Value, treasure[i].Weight = 0, 0
		}
	}

	resp := new(model.TreasureResponse)
	resp.StatusCode = 200
	resp.StatusMsg = "success"
	resp.Treasure = treasure
	resp.Capacity = capacity
	resp.Limit = limit
	c.JSON(http.StatusOK, resp)
}
