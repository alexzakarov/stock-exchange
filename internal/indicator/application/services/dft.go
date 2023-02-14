package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"

	ti "github.com/cinar/indicator"

	"math"
)

// IndicatorDft Indicator Calculate Func
func (u *service) IndicatorDft(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorDft begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var signal int
	var result []float64
	value0 := 0.0
	value1 := 0.0
	value2 := 0.0
	fisher := 0.0
	period := 34

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	highLine = ti.Max(period, closeLine)
	lowLine = ti.Min(period, closeLine)
	for i := 0; i < len(highLine); i++ {
		value2 = value0

		value0 = 0.66*((closeLine[i]-lowLine[i])/fx.FindMax(highLine[i]-lowLine[i], 0.001)-0.5) + 0.67*fx.CheckNaN(value2)
		if value0 > .99 {
			value1 = .999
		} else if value0 < -.99 {
			value1 = -.999
		} else {
			value1 = value0
		}

		fisher = .5*math.Log((1+value1)/fx.FindMax(1-value1, .001)) + .5*(fx.CheckNaN(fisher))
		if fisher > 0 {
			result = append(result, 1)
		} else if fisher < 0 {
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
