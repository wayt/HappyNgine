package redis

import (
	"errors"
	"github.com/wayt/happyngine/env"
	goredis "gopkg.in/redis.v3"
	"time"
)

var Client *goredis.Client

var CacheMiss = errors.New("redis: cache miss")

func init() {

	Client = goredis.NewClient(&goredis.Options{
		Addr:     env.Get("HAPPY_REDIS_HOST"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func Get(key string) (string, error) {

	res, err := Client.Get(key).Result()
	if err != nil {

		if err == goredis.Nil {
			return "", CacheMiss
		}

		return "", err
	}

	return res, nil
}

func Set(key, value string, expiration time.Duration) error {

	return Client.Set(key, value, expiration).Err()
}
