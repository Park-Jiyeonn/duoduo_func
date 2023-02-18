package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

func main() {
	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址
		Password: "",               // Redis密码，如果没有则留空
		DB:       0,                // Redis数据库索引（0表示默认数据库）
	})

	// 检查Redis连接是否正常
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	// 在Redis中设置一个键值对
	err = client.Set("T-ara", "Park.Jiyeon", 0).Err()
	if err != nil {
		panic(err)
	}

	// 从Redis中获取一个键值对
	val, err := client.Get("T-ara").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("T-ara:", val)

	// 关闭Redis客户端
	err = client.Close()
	if err != nil {
		panic(err)
	}
}
