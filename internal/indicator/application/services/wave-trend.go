package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorWaveTrend Indicator Calculate Func
func (u *service) IndicatorWaveTrend(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorWaveTrend begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var hlc3 []float64
	var abs []float64
	var ci []float64
	var signal int
	n1 := 10
	n2 := 21

	for _, values := range chartData {
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
		closeLine = append(closeLine, values.ClosePrice)
		hlc3 = append(hlc3, (values.HighPrice+values.LowPrice+values.ClosePrice)/3)
	}

	esa := fx.Ema(n1, hlc3)
	for i := 0; i < len(chartData); i++ {
		if hlc3[i]-esa[i] != 0 {
			abs = append(abs, math.Abs(hlc3[i]-esa[i]))
		} else {
			abs = append(abs, 1)
		}
	}

	d := fx.Ema(n1, abs)
	for i := 0; i < len(chartData); i++ {
		ci = append(ci, (hlc3[i]-esa[i])/(0.015*d[i]))
	}
	tci := fx.Ema(n2, ci)
	wt1 := tci
	wt2 := fx.Sma(4, wt1)
	cross := fx.Cross(wt1, wt2)

	signal = int(fx.GetByIndexN(cross, 2))

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
