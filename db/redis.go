package db

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"

	"helper/common/logger"
)

var redisConn *redis.Client
var once sync.Once

func GetRedisInstance() *redis.Client {
	once.Do(getRedis)
	return redisConn
}

func getRedis() {
	addr := viper.GetString("redis.server")
	if addr == "" {
		logger.Instance.Panic("get redis server form config empty")
	}
	logger.Instance.Debugf("pass: %v", viper.GetString("redis.pass"))
	c := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         addr,
		Password:     viper.GetString("redis.pass"),
		DB:           viper.GetInt("redis.db"),
		DialTimeout:  60 * time.Second,
		PoolSize:     1000,
		PoolTimeout:  2 * time.Minute,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	})
	_, err := c.Ping().Result()
	if err != nil {
		logger.Instance.Panicf("init redis err:" + err.Error())
	}
	redisConn = c
}
