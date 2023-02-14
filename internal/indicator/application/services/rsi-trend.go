package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorRsiTrend Indicator Calculate Func
func (u *service) IndicatorRsiTrend(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorRsiTrend begin to work")
	}

	var closeLine []float64
	var signal int
	var result []float64
	period := 14
	upperThreshold := 60.0
	lowerThreshold := 40.0

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	_, rsi := fx.RsiPeriod(period, closeLine)
	for i := 0; i < len(rsi); i++ {
		if rsi[i] > upperThreshold {
			result = append(result, 1)
		} else if rsi[i] < lowerThreshold {
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

	return responser, err
}
