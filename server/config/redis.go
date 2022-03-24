package config

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Rds struct{}

func Redis() Rds {
	return Rds{}
}

var client *redis.Client
var connectionErr error

func init() {
	redisAddr := GetEnv("REDIS_ADDR")
	redisUser := GetEnv("REDIS_USER")
	redisPassword := GetEnv("REDIS_PASSWORD")

	if IsTestMode() {
		redisAddr = "localhost:6379"
		redisUser = ""
		redisPassword = ""
	}

	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     redisAddr, // HOST:PORT
			Username: redisUser,
			Password: redisPassword,
		})
	}

	connectionErr = client.Ping(context.Background()).Err()

}

func (Rds) Client() (*redis.Client, error) {
	return client, connectionErr
}
