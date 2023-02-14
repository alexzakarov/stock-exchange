package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorSuperTrend Indicator Calculate Func
func (u *service) IndicatorSuperTrend(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorSuperTrend begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var hl21 []float64
	var up []float64
	var up1 []float64
	var down []float64
	var down1 []float64
	var trend []float64
	var signal int
	periods := 10
	multiplier := 3.0

	for _, values := range chartData {
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
		closeLine = append(closeLine, values.ClosePrice)
		hl21 = append(hl21, fx.Hl2(values.HighPrice, values.LowPrice))
	}

	_, atr2 := fx.Atr(periods, highLine, lowLine, closeLine)
	for i := 0; i < len(hl21); i++ {
		up = append(up, hl21[i]-(multiplier*atr2[i]))
		down = append(down, hl21[i]+(multiplier*atr2[i]))
	}

	for i := 1; i < len(up); i++ {
		up1 = append(up1, up[i-1])
		down1 = append(down1, down[i-1])
	}

	up = fx.SliceByIndex(up, 1, 0)
	down = fx.SliceByIndex(down, 1, 0)
	closeLine2 := fx.SliceByIndex(closeLine, 1, 0)
	for i := 2; i < len(down); i++ {
		if closeLine2[i-2] > up1[i-1] {
			up[i-1] = math.Max(up[i-1], up1[i-1])
		} else {
			up[i-1] = up[i-1]
		}
		if closeLine2[i-2] < down1[i-1] {
			down[i-1] = math.Min(down[i-1], down1[i-1])
		} else {
			down[i-1] = down[i-1]
		}
		up1[i] = up[i-1]
		down1[i] = down[i-1]
	}

	lastSignal := 0
	for i := 0; i < len(down1); i++ {
		if closeLine2[i] > down[i] && lastSignal != 1 {
			lastSignal = 1
			trend = append(trend, float64(lastSignal))
		} else if closeLine2[i] < up[i] && lastSignal != -1 {
			lastSignal = -1
			trend = append(trend, float64(lastSignal))

		} else {
			trend = append(trend, 0)
		}
	}

	signal = int(fx.GetByIndexN(trend, 2))

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
