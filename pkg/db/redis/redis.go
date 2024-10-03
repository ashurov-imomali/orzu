package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/config"
)

func NewRedisClient(conf config.RDb) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     conf.Url,
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.Database,
	})
	_, err := c.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return c, nil
}
