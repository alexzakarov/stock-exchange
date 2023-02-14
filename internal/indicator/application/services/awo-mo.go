package services

import (
	"context"
	ti "github.com/cinar/indicator"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorAwoMo Indicator Calculate Func
func (u *service) IndicatorAwoMo(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorAwoMo begin to work")
	}

	fastLength := 5
	slowLength := 34
	var hl2 []float64
	var ao []float64
	var signal int
	var result []float64

	for _, values := range chartData {
		hl2 = append(hl2, fx.Hl2(values.HighPrice, values.LowPrice))
	}

	fastSma := ti.Sma(fastLength, hl2)
	slowSma := ti.Sma(slowLength, hl2)
	for i := 0; i < len(fastSma); i++ {
		ao = append(ao, fastSma[i]-slowSma[i])
	}

	for i := 0; i < len(ao); i++ {
		if ao[i] >= 0 {
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
