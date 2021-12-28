package core

import (
	"context"
	"time"

	"git.myarena7.com/arena/hicourt/config"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

var client *redis.Client

func NewRedis() (*redis.Client, error) {

	conf := config.C().Redis

	client = redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         conf.Addr,
		DB:           conf.DB,
		Password:     conf.Pwd,
		PoolSize:     15,
		MinIdleConns: 8,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  3 * time.Second,
	})
	pong, err := client.Ping(context.TODO()).Result()
	if err == redis.Nil {
		log.Error().Err(err).Msg("redis异常")
		return nil, err
	} else if err != nil {
		log.Error().Err(err).Msg("连接redis失败")
		return nil, err
	} else {
		log.Info().Msgf("redis pong ,%v=", pong)
	}
	return client, nil
}
