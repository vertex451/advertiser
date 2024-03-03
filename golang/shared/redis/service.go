package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Service struct {
	rdb *redis.Client
}

func New() *Service {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		zap.L().Fatal("Error connecting to Redis", zap.Error(err))
	}

	return &Service{
		rdb: rdb,
	}
}

func (s *Service) Push(key, val string) error {
	// Push message onto Redis list
	ctx := context.Background()
	err := s.rdb.LPush(ctx, key, val).Err()
	if err != nil {
		zap.L().Error("Error pushing message onto Redis list", zap.Error(err))
	}

	return err
}

func (s *Service) Pop(key string) (string, error) {
	// Pop message from Redis list
	ctx := context.Background()
	msg, err := s.rdb.RPop(ctx, key).Result()
	if err != nil {
		zap.L().Error("Error popping message from Redis list", zap.Error(err))
		return "", err
	}
	zap.L().Info("Message from Redis list", zap.String("message", msg))

	return msg, nil
}
