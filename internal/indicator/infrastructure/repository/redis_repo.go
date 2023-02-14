package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"main/internal/indicator/domain/ports"
	"main/pkg/market_data/binance"
	"time"
)

// redisRepo Struct
type redisRepo struct {
	redisClient *redis.Client
}

// NewRedisRepo Indicator Domain redis repository constructor
func NewRedisRepo(redisClient *redis.Client) ports.IRedisRepository {
	return &redisRepo{redisClient: redisClient}
}

// GetCache Get cache by key
func (n *redisRepo) GetCache(ctx context.Context, key string) ([]binance.ChartData, error) {
	bytes, err := n.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "marketdataRedisRepo.SetCache.redisClient.GetCache")
	}

	var base []binance.ChartData
	if err = json.Unmarshal(bytes, &base); err != nil {
		return nil, errors.Wrap(err, "marketdataRedisRepo.SetCache.redisClient.GetCache.json.Unmarshal")
	}

	return base, nil
}

// SetCache Setting cache
func (n *redisRepo) SetCache(ctx context.Context, key string, seconds int, data any) error {
	indicatorsBytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "marketdataRedisRepo.SetCache.json.Marshal")
	}
	if err = n.redisClient.Set(ctx, key, indicatorsBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return errors.Wrap(err, "indicatorsRedisRepo.SetCache.redisClient.SetCache")
	}
	return nil
}

// DeleteCache Delete cache
func (n *redisRepo) DeleteCache(ctx context.Context, key string) error {
	if err := n.redisClient.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "marketdataRedisRepo.DeleteCache.redisClient.Del")
	}
	return nil
}
