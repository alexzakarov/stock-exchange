package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorRsi Indicator Calculate Func
func (u *service) IndicatorRsi(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorRsi begin to work")
	}

	var closeLine []float64
	var signal int
	var result []float64
	period := 14

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	_, rsi := fx.RsiPeriod(period, closeLine)

	for i := 0; i < len(rsi); i++ {
		if rsi[i] > 70 {
			result = append(result, 1)
		} else if rsi[i] < 30 {
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
