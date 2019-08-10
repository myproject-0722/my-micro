package gateway

import (
	conf "github.com/myproject-0722/my-micro/conf"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(
		&redis.Options{
			Addr: conf.RedisIP,
			DB:   2,
		},
	)

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Error(err)
		panic(err)
	}
}
