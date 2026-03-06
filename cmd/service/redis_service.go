package service

import (
	"chickchirick-auth/cmd/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type RedisDecorator struct {
	Client *redis.Client
}

func InitRedis(config config.RedisConfig) RedisDecorator {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Username: config.User,
		Password: config.Password,
	})

	redisClient := RedisDecorator{
		Client: client,
	}

	return redisClient
}

func (rd RedisDecorator) RedisClose() {
	err := rd.Client.Close()
	if err != nil {
		panic(fmt.Errorf("redis close error: %w", err))
	}
}
