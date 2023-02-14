package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorMbbt Indicator Calculate Func
func (u *service) IndicatorMbbt(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorMbbt begin to work")
	}

	var src []float64
	var signal int
	var result []float64
	fastMALength := 34
	slowMALength := 55
	refMALength := 5

	for _, values := range chartData {
		src = append(src, fx.Hl2(values.HighPrice, values.LowPrice))
	}

	signal = 0
	refMA := fx.Ema(refMALength, src)
	fastMA := fx.Ema(fastMALength, src)
	slowMA := fx.Ema(slowMALength, src)
	for i := 0; i < len(src); i++ {
		if refMA[i] > fastMA[i] && refMA[i] > slowMA[i] {
			result = append(result, 1)
		} else if refMA[i] < fastMA[i] && refMA[i] < slowMA[i] {
			result = append(result, -1)
		} else if refMA[i] > fastMA[i] && refMA[i] < slowMA[i] {
			result = append(result, 1)
		} else if refMA[i] < fastMA[i] && refMA[i] > slowMA[i] {
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
