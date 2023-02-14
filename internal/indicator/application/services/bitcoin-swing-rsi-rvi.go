package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorBitcoinSwingRsiRvi Indicator Calculate Func
func (u *service) IndicatorBitcoinSwingRsiRvi(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorBitcoinSwingRsiRvi begin to work")
	}

	var ohlc4 []float64
	var maxChange []float64
	var minChange []float64
	var rsi []float64
	var upperInput []float64
	var lowerInput []float64
	var rvi []float64
	var arr10 []float64
	var arr90 []float64
	var signal int
	var result []float64
	len1 := 7
	lenx := 50
	length := 10

	for _, values := range chartData {
		ohlc4 = append(ohlc4, (values.OpenPrice+values.HighPrice+values.LowPrice+values.ClosePrice)/4)
		arr10 = append(arr10, 10)
		arr90 = append(arr90, 90)
	}

	change := fx.Change(ohlc4)
	for i := 0; i < len(chartData); i++ {
		maxChange = append(maxChange, math.Max(change[i], 0))
		minChange = append(minChange, -math.Min(change[i], 0))
	}

	up := fx.Rma(len1, maxChange)
	down := fx.Rma(len1, minChange)
	for i := 0; i < len(chartData); i++ {
		if down[i] == 0 {
			rsi = append(rsi, 100)
		} else if up[i] == 0 {
			rsi = append(rsi, 0)
		} else {
			rsi = append(rsi, 100-(100/(1+up[i]/down[i])))
		}
	}

	stdev := fx.StDev(length, ohlc4)
	for i := 0; i < len(chartData); i++ {
		if change[i] <= 0 {
			upperInput = append(upperInput, 0)
		} else {
			upperInput = append(upperInput, stdev[i])
		}
		if change[i] > 0 {
			lowerInput = append(lowerInput, 0)
		} else {
			lowerInput = append(lowerInput, stdev[i])
		}
	}

	upper := fx.Ema(lenx, upperInput)
	lower := fx.Ema(lenx, lowerInput)
	for i := 0; i < len(chartData); i++ {
		rvi = append(rvi, upper[i]/(upper[i]+lower[i])*100)
	}

	crossOver10 := fx.CrossUp(rsi, arr10)
	crossUnder10 := fx.CrossDown(rsi, arr10)
	crossOver90 := fx.CrossUp(rsi, arr90)
	crossUnder90 := fx.CrossDown(rsi, arr90)
	for i := 0; i < len(chartData); i++ {
		if (crossOver10[i] == 1 || crossUnder10[i] == 1) && rvi[i] < 50 {
			result = append(result, -1)
		} else if (crossOver90[i] == 1 || crossUnder90[i] == 1) && rvi[i] > 50 {
			result = append(result, 1)
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
