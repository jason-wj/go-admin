package cache

import (
	"github.com/go-redis/redis/v7"
	"go-admin/common/core/sdk/runtime"
	"time"
)

// NewRedis redis模式
func NewRedis(client *redis.Client, options *redis.Options) (*Redis, error) {
	if client == nil {
		client = redis.NewClient(options)
	}
	r := &Redis{
		client: client,
	}
	err := r.connect()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Redis cache implement
type Redis struct {
	client *redis.Client
}

func (*Redis) String() string {
	return "redis"
}

// connect connect test
func (r *Redis) connect() error {
	var err error
	_, err = r.client.Ping().Result()
	return err
}

func (r *Redis) Exist(prefix, key string) bool {
	key = prefix + runtime.IntervalTenant + key
	v, _ := r.client.Exists(key).Result()
	if v != 1 {
		return false
	}
	return true
}

// Get from key
func (r *Redis) Get(prefix, key string) (string, error) {
	key = prefix + runtime.IntervalTenant + key
	return r.client.Get(key).Result()
}

// Set value with key and expire time
func (r *Redis) Set(prefix, key string, val interface{}, expire int) error {
	key = prefix + runtime.IntervalTenant + key
	return r.client.Set(key, val, time.Duration(expire)*time.Second).Err()
}

// Del delete key in redis
func (r *Redis) Del(prefix, key string) error {
	key = prefix + runtime.IntervalTenant + key
	return r.client.Del(key).Err()
}

func (r *Redis) HashSet(expire int, prefix, key string, values map[string]interface{}) error {
	key = prefix + runtime.IntervalTenant + key
	err := r.client.HSet(key, values).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(key, time.Duration(expire)*time.Second).Err()
}

// HashGet from key
func (r *Redis) HashGet(prefix, key, field string) (string, error) {
	key = prefix + runtime.IntervalTenant + key
	return r.client.HGet(key, field).Result()
}

func (r *Redis) HashGetAll(prefix, key string) (map[string]string, error) {
	key = prefix + runtime.IntervalTenant + key
	return r.client.HGetAll(key).Result()
}

// HashDel delete key in specify redis's hashtable
func (r *Redis) HashDel(prefix, key string, field string) error {
	key = prefix + runtime.IntervalTenant + key
	return r.client.HDel(key, field).Err()
}

// Increase
func (r *Redis) Increase(prefix, key string) error {
	key = prefix + runtime.IntervalTenant + key
	return r.client.Incr(key).Err()
}

func (r *Redis) Decrease(prefix, key string) error {
	key = prefix + runtime.IntervalTenant + key
	return r.client.Decr(key).Err()
}

// Set ttl
func (r *Redis) Expire(prefix, key string, expire int) error {
	key = prefix + runtime.IntervalTenant + key
	return r.client.Expire(key, time.Duration(expire)*time.Second).Err()
}

// GetClient 暴露原生client
func (r *Redis) GetClient() *redis.Client {
	return r.client
}
