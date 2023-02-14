package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorMultiTrend Indicator Calculate Func
func (u *service) IndicatorMultiTrend(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorMultiTrend begin to work")
	}

	var closeLine []float64
	var signal int
	var result []float64
	period := 200

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	ema := fx.Ema(period, closeLine)
	for i := 1; i < len(chartData); i++ {
		if ema[i] > ema[i-1] {
			result = append(result, 1)
		} else if ema[i] < ema[i-1] {
			result = append(result, -1)
		} else {
			result = append(result, 0)
		}
	}

	signal = int(fx.GetByIndexN(result, 2))

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
