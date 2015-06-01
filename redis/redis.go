package redis

import (
	"github.com/wayt/happyngine/env"
	goredis "gopkg.in/redis.v3"
	"time"
)

var Client *goredis.Client

func init() {

	Client = goredis.NewClient(&goredis.Options{
		Addr:     env.Get("REDIS_PORT_6379_TCP_ADDR") + ":" + env.Get("REDIS_PORT_6379_TCP_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

type Item struct {
	Key        string
	Value      string
	Expiration time.Duration
	Missed     bool
}

func Get(key string) (*Item, error) {

	item := &Item{
		Key:    key,
		Missed: false,
	}

	var err error
	item.Value, err = Client.Get(key).Result()
	if err != nil {

		if err == goredis.Nil {
			item.Missed = true
			return item, nil
		}

		return nil, err
	}

	return item, nil
}

func Set(item *Item) error {

	return Client.Set(item.Key, item.Value, item.Expiration).Err()
}

func Delete(key string) error {
	return Client.Del(key).Err()
}
