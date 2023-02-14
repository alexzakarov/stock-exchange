package services

import (
	"context"
	ti "github.com/cinar/indicator"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	"math"
)

// IndicatorScc Calculate Func
func (u *service) IndicatorScc(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorScc begin to work")
	}

	var highLine []float64
	var lowLine []float64
	var closeLine []float64
	var hl2 []float64
	var trendUp []float64
	var trendDown []float64
	var trend []float64
	var up []float64
	var down []float64
	var signal int
	period1 := 9
	factor1 := 4.0

	for _, values := range chartData {
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
		closeLine = append(closeLine, values.ClosePrice)
		hl2 = append(hl2, (values.HighPrice+values.LowPrice)/2)
	}

	for index, value := range chartData {
		highLine = append(highLine, value.HighPrice)
		lowLine = append(lowLine, value.LowPrice)
		closeLine = append(closeLine, value.ClosePrice)

		_, atrValue := ti.Atr(period1, highLine, lowLine, closeLine)
		up = append(up, hl2[index]-(factor1*atrValue[index]))
		down = append(down, hl2[index]+(factor1*atrValue[index]))
	}

	for i := 1; i < len(up); i++ {
		if closeLine[i] > trendUp[i-1] {
			trendUp = append(trendUp, math.Max(up[i], trendUp[i-1]))
		} else {
			trendUp = append(trendUp, up[i])
		}

		if closeLine[i] < trendDown[i-1] {
			trendDown = append(trendDown, math.Min(down[i], trendDown[i-1]))
		} else {
			trendDown = append(trendDown, down[i])
		}

		if closeLine[i] > trendUp[i-1] {
			trend = append(trend, 1)
		} else if closeLine[i] < trendDown[i-1] {
			trend = append(trend, -1)
		} else {
			trend = append(trend, trend[i-1])
		}
	}

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
