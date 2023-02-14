package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"

	ti "github.com/cinar/indicator"
)

// IndicatorTutci Indicator Calculate Func
func (u *service) IndicatorTutci(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorTutci begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var highUpCond []bool
	var lowDownCond []bool
	var K1 []float64
	var K2 []float64
	var K3 []float64
	var K4 []float64
	var buySignal []bool
	var sellSignal []bool
	var buyExit []bool
	var sellExit []bool
	var upper1 []float64
	var lower1 []float64
	var sup1 []float64
	var sdown1 []float64
	var signal int
	var result []float64
	period1 := 20
	period2 := 10

	lowDownCond = append(lowDownCond, false)
	highUpCond = append(highUpCond, false)
	for _, values := range chartData {
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
		closeLine = append(closeLine, values.ClosePrice)
	}

	lower := ti.Min(period1, lowLine)
	upper := ti.Max(period1, highLine)
	down := ti.Min(period1, lowLine)
	up := ti.Max(period1, highLine)
	sDown := ti.Min(period2, lowLine)
	sUp := ti.Max(period2, highLine)
	for i := 1; i < len(chartData); i++ {
		lowDownCond = append(lowDownCond, lowLine[i] <= down[i-1])
		highUpCond = append(highUpCond, highLine[i] >= up[i-1])
	}

	lowDownCondBarSince := fx.BarSince(lowDownCond)
	highUpCondBarSince := fx.BarSince(highUpCond)
	for i := 0; i < len(chartData); i++ {
		if highUpCondBarSince[i] <= lowDownCondBarSince[i] {
			K1 = append(K1, down[i])
		} else {
			K1 = append(K1, up[i])
		}
		if highUpCondBarSince[i] <= lowDownCondBarSince[i] {
			K2 = append(K2, sDown[i])
		} else {
			K2 = append(K2, sUp[i])
		}
		if closeLine[i] > K1[i] {
			K3 = append(K3, down[i])
		} else {
			K3 = append(K3, 0)
		}
		if closeLine[i] < K1[i] {
			K4 = append(K4, up[i])
		} else {
			K4 = append(K4, 0)
		}
	}

	upper1 = append([]float64{0}, upper...)
	upper1 = upper1[:len(upper1)-1]
	lower1 = append([]float64{0}, lower...)
	lower1 = lower1[:len(lower1)-1]
	sup1 = append([]float64{0}, sUp...)
	sup1 = sup1[:len(sup1)-1]
	sdown1 = append([]float64{0}, sDown...)
	sdown1 = sdown1[:len(sdown1)-1]

	crossHighUpper := fx.CrossUp(highLine, upper1)
	crossLowerLow := fx.CrossUp(lower1, lowLine)
	crossSdownLow := fx.CrossUp(sdown1, lowLine)
	crossHighSup := fx.CrossUp(highLine, sup1)
	for i := 1; i < len(chartData); i++ {
		if highLine[i] == upper[i-1] || crossHighUpper[i] == 1 {
			buySignal = append(buySignal, true)
		} else {
			buySignal = append(buySignal, false)
		}

		if lowLine[i] == lower[i-1] || crossLowerLow[i] == 1 {
			sellSignal = append(sellSignal, true)
		} else {
			sellSignal = append(sellSignal, false)
		}

		if lowLine[i] == sDown[i-1] || crossSdownLow[i] == 1 {
			buyExit = append(buyExit, true)
		} else {
			buyExit = append(buyExit, false)
		}

		if highLine[i] == sUp[i-1] || crossHighSup[i] == 1 {
			sellExit = append(sellExit, true)
		} else {
			sellExit = append(sellExit, false)
		}
	}

	O1 := fx.BarSince(buySignal)
	O2 := fx.BarSince(sellSignal)
	O3 := fx.BarSince(buyExit)
	O4 := fx.BarSince(sellExit)

	for i := 1; i < len(O1); i++ {
		if buySignal[i] && O3[i] < O1[i-1] {
			result = append(result, 1)
		} else if sellSignal[i] && O4[i] < O2[i-1] {
			result = append(result, -1)
		} else {
			result = append(result, 0)
		}
	}

	signal = int(fx.GetByIndexN(result, 2))

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-2].CloseTime,
	}

	return
}
