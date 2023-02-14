package ports

import (
	"context"
	ent "main/internal/indicator/domain/entities"
)

// IPostgresqlRepository Indicator Domain postgresql interface
type IPostgresqlRepository interface {
	WriteResult(context.Context, *ent.IndicatorCalcResponse) error
	ReadIndicatorIndexByInterval(context.Context, string) ([]ent.IndicatorIndex, error)
	ChangeLockStatusByRow(context.Context, int64) error
	ReleaseAllLocks(context.Context) error
	UpdateIndicatorResult(context.Context, ent.IndicatorIndex, ent.IndicatorCalcResponse) error
	SaveBoosterStatistics(context.Context, string, int64, string, ent.BoosterStatisticsResponse) error
}
