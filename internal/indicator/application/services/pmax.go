package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorPMax Indicator Calculate Func
func (u *service) IndicatorPMax(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorPMax begin to work")
	}
	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var hl2 []float64
	var longStop []float64
	var shortStop []float64
	var longStopPrev []float64
	var shortStopPrev []float64
	var dir []float64
	var pmax []float64
	var signal int
	var result []float64

	period := 10.0
	mult := 3.0
	length := 10.0

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
		hl2 = append(hl2, (values.HighPrice+values.LowPrice)/2)
		dir = append(dir, 1)
		longStopPrev = append(longStopPrev, 0)
		shortStopPrev = append(shortStopPrev, 0)

	}

	_, atr := fx.Atr(int(period), highLine, lowLine, closeLine)

	mavg := fx.Ema(int(length), hl2)

	for i := 0; i < len(chartData); i++ {
		longStop = append(longStop, mavg[i]-mult*atr[i])
		shortStop = append(shortStop, mavg[i]+mult*atr[i])
	}

	for i := 1; i < len(chartData); i++ {
		longStopPrev[i] = longStop[i-1]
		shortStopPrev[i] = shortStop[i-1]
		if mavg[i] > longStopPrev[i] {
			longStop[i] = math.Max(longStop[i], longStopPrev[i])
		}
		if mavg[i] < shortStopPrev[i] {
			shortStop[i] = math.Min(shortStop[i], shortStopPrev[i])
		}
	}

	for i := 1; i < len(chartData); i++ {
		dir[i] = dir[i-1]
		if dir[i] == -1 && mavg[i] > shortStopPrev[i] {
			dir[i] = 1
		} else if dir[i] == 1 && mavg[i] < longStopPrev[i] {
			dir[i] = -1
		}
	}

	for i := 0; i < len(chartData); i++ {
		if dir[i] == 1 {
			pmax = append(pmax, longStop[i])
		} else {
			pmax = append(pmax, shortStop[i])
		}
	}

	buySignalK := fx.CrossUp(mavg, pmax)
	sellSignalK := fx.CrossDown(mavg, pmax)

	for i := 0; i < len(chartData); i++ {
		if buySignalK[i] == 1 {
			result = append(result, 1)
		} else if sellSignalK[i] == 1 {
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
	result = nil
	return responser, err
}
