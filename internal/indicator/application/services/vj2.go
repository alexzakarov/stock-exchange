package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorVj2 Indicator Calculate Func
func (u *service) IndicatorVj2(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorVj2 begin to work")
	}

	const factor float64 = 3
	const period int = 14
	var highLine []float64
	var lowLine []float64
	var closeLine []float64
	var hl2 []float64
	var up []float64
	var down []float64
	var trend []float64
	var trendUp []float64
	var trendDown []float64
	var signal int

	for _, values := range chartData {
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
		closeLine = append(closeLine, values.ClosePrice)
		hl2 = append(hl2, fx.Hl2(values.HighPrice, values.LowPrice))
	}

	_, atr := fx.Atr(period, highLine, lowLine, closeLine)
	for i := 0; i < len(chartData); i++ {
		up = append(up, hl2[i]-(factor*atr[i]))
		down = append(down, hl2[i]+(factor*atr[i]))
	}

	trendUp = append(trendUp, up[0])
	trendDown = append(trendDown, down[0])
	trend = append(trend, 0)
	for i := 1; i < len(chartData); i++ {
		if closeLine[i-1] > trendUp[i-1] {
			trendUp = append(trendUp, math.Max(trendUp[i-1], up[i]))
		} else {
			trendUp = append(trendUp, up[i])
		}

		if closeLine[i-1] < trendDown[i-1] {
			trendDown = append(trendDown, math.Min(trendDown[i-1], down[i]))
		} else {
			trendDown = append(trendDown, down[i])
		}
	}

	for i := 1; i < len(chartData); i++ {
		if closeLine[i] > trendDown[i-1] {
			trend = append(trend, 1)
		} else if closeLine[i] < trendUp[i-1] {
			trend = append(trend, -1)
		} else {
			trend = append(trend, trend[i-1])
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
