package datasource

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	lib_error "github.com/ossrs/go-oryx-lib/errors"
	"os"
	"strconv"
)

type RedisDatasource struct {
	rdb *redis.Client
}

func envRedisDatabase() string {
	return os.Getenv("REDIS_DATABASE")
}

func envRedisHost() string {
	return os.Getenv("REDIS_HOST")
}

func envRedisPort() string {
	return os.Getenv("REDIS_PORT")
}
func envRedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}

func NewRedisDatasource() (Datasource, error) {
	redisDatabase, err := strconv.Atoi(envRedisDatabase())
	if err != nil {
		return nil, lib_error.Wrapf(err, "invalid REDIS_DATABASE %v", envRedisDatabase())
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", envRedisHost(), envRedisPort()),
		Password: envRedisPassword(),
		DB:       redisDatabase,
	})
	return &RedisDatasource{rdb: rdb}, nil
}

func (r *RedisDatasource) Get(ctx context.Context, key, field string) (string, error) {
	var err error
	var data string
	if field == "" {
		data, err = r.rdb.Get(ctx, key).Result()
	} else {
		data, err = r.rdb.HGet(ctx, key, field).Result()
	}
	// Nil reply returned by Redis when key does not exist.
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}
	return data, nil
}

func (r *RedisDatasource) Set(ctx context.Context, key, field string, value string) error {
	var err error
	if field == "" {
		err = r.rdb.Set(ctx, key, field, 0).Err()
	} else {
		err = r.rdb.HSet(ctx, key, field, value).Err()
	}

	// Nil reply returned by Redis when key does not exist.
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}

func (r *RedisDatasource) Delete(ctx context.Context, key, field string) error {
	err := r.rdb.HDel(ctx, key, field).Err()
	// Nil reply returned by Redis when key does not exist.
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}

func (r *RedisDatasource) DeleteAll(ctx context.Context, key string) error {
	err := r.rdb.Del(ctx, key).Err()
	// Nil reply returned by Redis when key does not exist.
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}

func (r *RedisDatasource) Select(ctx context.Context, key string, options *SelectOptions) ([]string, error) {
	if options.Filed == "" {
		options.Filed = "*"
	}
	results, _, err := r.rdb.HScan(ctx, key, 0, options.Filed, options.Count).Result()
	// Nil reply returned by Redis when key does not exist.
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	return results, nil
}

func (r *RedisDatasource) SelectAll(ctx context.Context, key string) (map[string]string, error) {
	keys, err := r.rdb.HGetAll(ctx, key).Result()
	// Nil reply returned by Redis when key does not exist.
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	return keys, nil
}

func (r *RedisDatasource) Count(ctx context.Context, key string) (int64, error) {
	r1, err := r.rdb.HLen(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}
	return r1, nil
}

// Incr 当前元素原子操作
func (r *RedisDatasource) Incr(ctx context.Context, key, field string, value int64) error {
	err := r.rdb.HIncrBy(ctx, key, field, value).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return nil
}
