package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"

	ti "github.com/cinar/indicator"
)

// IndicatorMacd Indicator Calculate Func
func (u *service) IndicatorMacd(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorMacd begin to work")
	}

	var closeLine []float64
	var histogram []float64
	var signal int
	var result []float64
	period := 9

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	macd, _ := ti.Macd(closeLine)
	signalEma := fx.Ema(period, macd)
	for i := 0; i < len(chartData); i++ {
		histogram = append(histogram, macd[i]-signalEma[i])
	}

	for i := 0; i < len(chartData); i++ {
		if histogram[i] > 0 {
			result = append(result, 1)
		} else if histogram[i] < 0 {
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
