package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorPsar Indicator Calculate Func
func (u *service) IndicatorPsar(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorPsar begin to work")
	}

	var highLine []float64
	var lowLine []float64
	var closeLine []float64
	var trendSignal []fx.Trend
	var trend []fx.Trend
	var signal int

	for _, values := range chartData {
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
		closeLine = append(closeLine, values.ClosePrice)
	}

	_, trend = fx.ParabolicSar(highLine, lowLine, closeLine, 0.00133, 0.22)
	for i := 1; i < len(trend); i++ {
		if trend[i] != trend[i-1] {
			trendSignal = append(trendSignal, trend[i])
		} else {
			trendSignal = append(trendSignal, 0)
		}
	}

	signal = int(trendSignal[len(trendSignal)-2])

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return responser, err
}
