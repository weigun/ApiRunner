package dao

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type Cache interface {
	// get cached value by key.
	Get(key string) string
	// GetMulti is a batch version of Get.
	// GetMulti(keys []string) []interface{}
	// set cached value with key and expire time.
	Put(key string, val interface{}, timeout time.Duration) error
	MPut(keys, vals []interface{}) error
	// delete cached value by key.
	Delete(key string) error
	// increase cached int value by key, as a counter.
	Incr(key string) error
	// decrease cached int value by key, as a counter.
	Decr(key string) error
	// check if cached value exists or not.
	IsExist(key string) bool
	Keys(pattern string) []string
	// clear all cache.
	// ClearAll() error
}

var redisCachePtr *redisCache

var once sync.Once

func newCache() {
	once.Do(func() {
		//TODO 配置化
		client := redis.NewClient(&redis.Options{
			Addr:     "192.168.128.134:6379",
			Password: os.Getenv("GATE"),
			DB:       0,
		})
		_, err := client.Ping().Result()
		if err != nil {
			panic(err)
		}
		redisCachePtr = &redisCache{client}
	})
	// return redisCachePtr

}

func GetCache() *redisCache {
	return redisCachePtr
}

var (
	// DefaultCacheKey the collection name of redis for cache adapter.
	DefaultCacheKey = "ApiRunnerCache"
)

type redisCache struct {
	Client *redis.Client
}

func (this *redisCache) Get(key string) string {
	key = fmt.Sprintf(`%s:%s`, DefaultCacheKey, key)
	if v, err := this.Client.Get(key).Result(); err == nil {
		return v
	}
	return ``
}

func (this *redisCache) Put(key string, val interface{}, timeout time.Duration) error {
	key = fmt.Sprintf(`%s:%s`, DefaultCacheKey, key)
	_, err := this.Client.Set(key, val, timeout/time.Second).Result()
	return err
}

func (this *redisCache) MPut(keys []interface{}, vals []interface{}) error {
	// format key
	if len(keys) != len(vals) {
		return errors.New(`Key-value pairs are not equal`)
	}
	innerKeys := make([]interface{}, len(keys), cap(keys))
	copy(innerKeys, keys)
	for i, k := range innerKeys {
		innerKeys[i] = fmt.Sprintf(`%s:%s`, DefaultCacheKey, k.(string))
	}
	//combine k-v pairs
	kvPairs := []interface{}{}
	for i := 0; i < len(innerKeys); i++ {
		kvPairs = append(kvPairs, innerKeys[i])
		kvPairs = append(kvPairs, vals[i])
	}
	_, err := this.Client.MSet(kvPairs...).Result()
	return err
}

func (this *redisCache) Delete(key string) error {
	key = fmt.Sprintf(`%s:%s`, DefaultCacheKey, key)
	_, err := this.Client.Del(key).Result()
	return err
}

func (this *redisCache) Incr(key string) error {
	key = fmt.Sprintf(`%s:%s`, DefaultCacheKey, key)
	_, err := this.Client.Incr(key).Result()
	return err
}

func (this *redisCache) Decr(key string) error {
	key = fmt.Sprintf(`%s:%s`, DefaultCacheKey, key)
	_, err := this.Client.Decr(key).Result()
	return err
}

func (this *redisCache) IsExist(key string) bool {
	key = fmt.Sprintf(`%s:%s`, DefaultCacheKey, key)
	if _, err := this.Client.Exists(key).Result(); err == nil {
		return true
	}
	return false
}

func (this *redisCache) Keys(pattern string) []string {
	key := fmt.Sprintf(`%s:%s`, DefaultCacheKey, pattern)
	if v, err := this.Client.Keys(key).Result(); err == nil {
		return v
	}
	return []string{}
}

func init() {
	log.Println(`redis dao init`)
	newCache()
}
