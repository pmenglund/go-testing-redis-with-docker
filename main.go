package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func getNext(redis *redis.Client) int64 {
	return redis.Incr("key").Val()
}

func main() {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	val := getNext(client)
	fmt.Printf("value is %d\n", val)
}
