package redis_manager

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jimersylee/iris-seed/config"
	"github.com/sirupsen/logrus"
)

// client
var client *redis.Client

func InitRedisManager() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Addr,
		Password: config.Conf.Redis.Password, // no password set
		DB:       0,                          // use default DB
	})

	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	client = redisClient
	logrus.Info(client)
}

// GetClient /**
func GetClient() *redis.Client {
	if client == nil {
		InitRedisManager()
	}
	return client
}
