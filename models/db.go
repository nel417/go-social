package models

import (
	"github.com/go-redis/redis"
)

//global redis client
var client *redis.Client

func Init() {
	//redis client and host
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
