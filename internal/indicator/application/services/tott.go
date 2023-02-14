package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorTOTT Indicator Calculate Func
func (u *service) IndicatorTOTT(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorTOTT begin to work")
	}

	var closeLine []float64
	var diff []float64
	var longStop []float64
	var longStopPrev []float64
	var shortStop []float64
	var shortStopPrev []float64
	var dir []float64
	var MT []float64
	var OTT []float64
	var OTTup []float64
	var OTTup2 []float64
	var OTTdown []float64
	var OTTdown2 []float64
	var buySignalK1 []float64
	var sellSignalK1 []float64
	var signal int
	var result []float64

	period := 40
	percent := 1.0
	coeff := 0.001

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		dir = append(dir, 1)
	}

	mav := fx.TOTTMA(closeLine, period)

	for i := 0; i < len(chartData); i++ {
		diff = append(diff, mav[i]*percent*0.01)
		longStop = append(longStop, mav[i]-diff[i])
		shortStop = append(shortStop, mav[i]+diff[i])
	}
	longStopPrev = append(longStopPrev, 0)
	shortStopPrev = append(shortStopPrev, 0)

	for i := 1; i < len(chartData); i++ {
		longStopPrev = append(longStopPrev, longStop[i-1])
		shortStopPrev = append(shortStopPrev, shortStop[i-1])
		if mav[i] > longStopPrev[i] {
			longStop[i] = math.Max(longStop[i], longStopPrev[i])
		}
		if mav[i] < shortStopPrev[i] {
			shortStop[i] = math.Min(shortStop[i], shortStopPrev[i])
		}
	}

	for i := 1; i < len(chartData); i++ {
		dir[i] = dir[i-1]
		if dir[i] == -1 && mav[i] > shortStopPrev[i] {
			dir[i] = 1
		} else if dir[i] == 1 && mav[i] < longStopPrev[i] {
			dir[i] = -1
		}
	}

	for i := 0; i < len(chartData); i++ {
		if dir[i] == 1 {
			MT = append(MT, longStop[i])
		} else {
			MT = append(MT, shortStop[i])
		}
		if mav[i] > MT[i] {
			OTT = append(OTT, MT[i]*(200+percent)/200)
		} else {
			OTT = append(OTT, MT[i]*(200-percent)/200)
		}
		OTTup = append(OTTup, OTT[i]*(1+coeff))
		OTTdown = append(OTTdown, OTT[i]*(1-coeff))
	}

	OTTup2 = append(OTTup2, 0, 0)
	OTTdown2 = append(OTTdown2, 0, 0)
	for i := 2; i < len(chartData); i++ {
		OTTup2 = append(OTTup2, OTTup[i-2])
		OTTdown2 = append(OTTdown2, OTTdown[i-2])
	}

	buySignalK := fx.CrossUp(mav, OTTup2)
	sellSignalK := fx.CrossDown(mav, OTTdown2)

	buySignalK1 = append(buySignalK1, 0)
	sellSignalK1 = append(sellSignalK1, 0)
	for i := 1; i < len(chartData); i++ {
		buySignalK1 = append(buySignalK1, buySignalK[i-1])
		sellSignalK1 = append(sellSignalK1, sellSignalK[i-1])
	}

	K1 := fx.BarSince2(buySignalK)
	K2 := fx.BarSince2(sellSignalK)
	O1 := fx.BarSince2(buySignalK1)
	O2 := fx.BarSince2(sellSignalK1)

	for i := 0; i < len(chartData); i++ {
		if buySignalK[i] == 1 && O1[i] > K2[i] {
			result = append(result, 1)
		} else if sellSignalK[i] == 1 && O2[i] > K1[i] {
			result = append(result, -1)
		} else {
			result = append(result, 0)
		}
	}

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}
	result = nil
	return responser, err
}
