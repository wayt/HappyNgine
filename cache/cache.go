package cache

import (
	gocache "github.com/pmylund/go-cache"
	"github.com/wayt/happyngine/env"
	"time"
)

var Store *gocache.Cache

func init() {

	defExpire := env.GetInt("HAPPYNGINE_DEFAULT_CACHE_EXPIRATION")
	if defExpire <= 0 {
		defExpire = 5
	}

	cleanInterval := defExpire / 2
	if cleanInterval < 1 {
		cleanInterval = 1
	}

	Store = gocache.New(time.Duration(defExpire)*time.Second, time.Duration(cleanInterval)*time.Second)
}

const (
	// For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
	// For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to New() or
	// NewFrom() when the cache was created (e.g. 5 minutes.)
	DefaultExpiration time.Duration = 0
)

func Set(k string, x interface{}, d time.Duration) {
	Store.Set(k, x, d)
}

func Get(k string) (interface{}, bool) {
	return Store.Get(k)
}

func Delete(k string) {
	Store.Delete(k)
}
