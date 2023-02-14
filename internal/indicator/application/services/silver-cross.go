package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorSilverCross Indicator Calculate Func
func (u *service) IndicatorSilverCross(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorSilverCross begin to work")
	}

	var closeLine []float64
	var signal int
	var result []float64
	period1 := 20
	period2 := 50

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	ma20 := fx.Ema(period1, closeLine)
	ma50 := fx.Ema(period2, closeLine)
	for i := 0; i < len(chartData); i++ {
		if ma20[i] > ma50[i] {
			result = append(result, 1)
		} else {
			result = append(result, -1)
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
