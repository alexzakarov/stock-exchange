package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorCMSlingShotSystem Indicator Calculate Func
func (u *service) IndicatorCMSlingShotSystem(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorCMSlingShotSystem begin to work")
	}

	var closeLine []float64
	var upTrend []float64
	var downTrend []float64
	var signal int
	var result []float64
	emaSlowPeriod := 62
	emaFastPeriod := 38

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	emaSlow := fx.Ema(emaSlowPeriod, closeLine)
	emaFast := fx.Ema(emaFastPeriod, closeLine)
	for i := 0; i < len(chartData); i++ {
		if emaFast[i] > emaSlow[i] {
			upTrend = append(upTrend, 1)
		} else {
			upTrend = append(upTrend, 0)
		}
		if emaFast[i] < emaSlow[i] {
			downTrend = append(downTrend, 1)
		} else {
			downTrend = append(downTrend, 0)
		}
	}

	for i := 0; i < len(upTrend); i++ {
		if upTrend[i] == 1 {
			result = append(result, 1)
		} else if downTrend[i] == 1 {
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
