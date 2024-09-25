package redis

import (
	"fmt"
	"project_bluebell/settings"

	"github.com/go-redis/redis"
)

var Client *redis.Client
var Nil = redis.Nil

func Init(cfg *settings.RedisConfig) (err error) {

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.Poolsize,
	})

	_, err = Client.Ping().Result()
	return err
}

func Close() {
	_ = Client.Close()
}
