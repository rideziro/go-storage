package cache

import (
	"context"
	"errors"
	"fmt"
	redisClient "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

var (
	ErrRetrieve  = errors.New("unable to retrieve")
	ErrStore     = errors.New("unable to store")
	ErrSetExpiry = errors.New("unable to set expiry")
)

type Redis struct {
	client *redisClient.Client
}

func NewRedis(uri, port, password string) Redis {
	return Redis{
		redisClient.NewClient(&redisClient.Options{
			Addr:     fmt.Sprintf("%s:%s", uri, port),
			Password: password,
		}),
	}
}
func NewDefaultRedis() Redis {
	uri := viper.GetString("REDIS_URI")
	port := viper.GetString("REDIS_PORT")
	password := viper.GetString("REDIS_PASSWORD")
	return NewRedis(uri, port, password)
}

func (r *Redis) Store(ctx context.Context, values map[string]string) error {
	err := r.client.MSet(ctx, values).Err()
	if err != nil {
		return ErrStore
	}
	return nil
}

func (r *Redis) StoreWithExpire(ctx context.Context, values map[string]string, duration time.Duration) error {
	err := r.Store(ctx, values)
	if err != nil {
		return err
	}
	for key, _ := range values {
		err := r.client.Expire(ctx, key, duration).Err()
		if err != nil {
			return ErrSetExpiry
		}
	}

	return nil
}

func (r *Redis) Get(ctx context.Context, keys ...string) (map[string]string, error) {
	result, err := r.client.MGet(ctx, keys...).Result()

	if err != nil {
		return nil, ErrRetrieve
	}

	data := map[string]string{}
	for i, key := range keys {
		s, ok := result[i].(string)
		if ok {
			data[key] = s
		} else {
			data[key] = ""
		}
	}
	return data, nil
}
