package redis

import (
	"github.com/go-redis/redis"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/config"
	"go.uber.org/zap"
)

type RedisClient interface {
	Connect() (*redis.Client, error)
	Disconnect()
}

type redisClient struct {
	service *config.Service
	client  *redis.Client
	logger  *zap.Logger
}

func NewRedisClient(service *config.Service, logger *zap.Logger) RedisClient {
	return &redisClient{
		service: service,
		logger:  logger,
	}
}

func (c *redisClient) Connect() (*redis.Client, error) {
	redisDb := redis.NewClient(&redis.Options{
		Addr: c.service.GetConfig().RedisConnect,
	})
	_, err := redisDb.Ping().Result()
	if err != nil {
		return nil, err
	}
	return redisDb, nil
}

func (c *redisClient) Disconnect() {
	err := c.client.Close()
	if err != nil {
		c.logger.Error("redis disconnect error")
	}
}
