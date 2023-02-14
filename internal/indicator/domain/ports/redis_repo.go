package ports

import (
	"context"
	"main/pkg/market_data/binance"
)

// IRedisRepository Indicator Domain redis interface
type IRedisRepository interface {
	GetCache(ctx context.Context, key string) ([]binance.ChartData, error)
	SetCache(ctx context.Context, key string, seconds int, data any) error
	DeleteCache(ctx context.Context, key string) error
}
