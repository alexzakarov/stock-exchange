package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorSuperTrendFixed2 Indicator Calculate Func
func (u *service) IndicatorSuperTrendFixed2(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorSuperTrendFixed2 begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var hlc3 []float64
	var hl2 []float64
	var up []float64
	var down []float64
	var up2 []float64
	var down2 []float64
	var trendUp2 []float64
	var trendDown2 []float64
	var trend2 []float64
	var tsl2 []float64
	var trendUp []float64
	var trendDown []float64
	var trend []float64
	var tsl []float64
	var longcondition []float64
	var shortcondition []float64
	var longcondition2 []float64
	var shortcondition2 []float64
	var vortexBuyPermit []float64
	var vortexSellPermit []float64
	var longEntry []float64
	var shortEntry []float64
	var signal int
	var result []float64
	factor := 3
	pd := 14
	factor2 := 2
	pd2 := 10
	vortexBuyPeriod := 56
	vortexSellPeriod := 10

	for _, values := range chartData {
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
		closeLine = append(closeLine, values.ClosePrice)
		hlc3 = append(hlc3, (values.HighPrice+values.LowPrice+values.ClosePrice)/3)
		hl2 = append(hl2, fx.Hl2(values.HighPrice, values.LowPrice))
	}

	_, atrPd2 := fx.Atr(pd2, highLine, lowLine, closeLine)
	for i := 0; i < len(chartData); i++ {
		up2 = append(up2, hlc3[i]-float64(factor2)*atrPd2[i])
		down2 = append(down2, hlc3[i]+float64(factor2)*atrPd2[i])
	}

	trendUp2 = append(trendUp2, up2[0])
	trendDown2 = append(trendDown2, down2[0])
	trend2 = append(trend2, 0)
	tsl2 = append(tsl2, 0)
	for i := 1; i < len(chartData); i++ {
		if closeLine[i-1] > trendUp2[i-1] {
			trendUp2 = append(trendUp2, math.Max(up2[i], trendUp2[i-1]))
		} else {
			trendUp2 = append(trendUp2, up2[i])
		}

		if closeLine[i-1] < trendDown2[i-1] {
			trendDown2 = append(trendDown2, math.Min(down2[i], trendDown2[i-1]))
		} else {
			trendDown2 = append(trendDown2, down2[i])
		}

		if closeLine[i-1] > trendDown2[i-1] {
			trend2 = append(trend2, 1)
		} else if closeLine[i-1] < trendUp2[i-1] {
			trend2 = append(trend2, -1)
		} else {
			trend2 = append(trend2, trend2[i-1])
		}

		if trend2[i] == 1 {
			tsl2 = append(tsl2, trendUp2[i])
		} else {
			tsl2 = append(tsl2, trendDown2[i])
		}
	}

	for i := 0; i < len(chartData); i++ {
		if closeLine[i] > tsl2[i] {
			longcondition = append(longcondition, 1)
		} else {
			longcondition = append(longcondition, 0)
		}

		if closeLine[i] < tsl2[i] {
			shortcondition = append(shortcondition, 1)
		} else {
			shortcondition = append(shortcondition, 0)
		}
	}

	_, atrPd := fx.Atr(pd, highLine, lowLine, closeLine)
	for i := 0; i < len(chartData); i++ {
		up = append(up, hl2[i]-float64(factor)*atrPd[i])
		down = append(down, hl2[i]+float64(factor)*atrPd[i])
	}

	trendUp = append(trendUp, up2[0])
	trendDown = append(trendDown, down2[0])
	trend = append(trend, 0)
	tsl = append(tsl, 0)
	for i := 1; i < len(chartData); i++ {
		if closeLine[i-1] > trendUp[i-1] {
			trendUp = append(trendUp, math.Max(up[i], trendUp[i-1]))
		} else {
			trendUp = append(trendUp, up[i])
		}

		if closeLine[i-1] < trendDown[i-1] {
			trendDown = append(trendDown, math.Min(down[i], trendDown[i-1]))
		} else {
			trendDown = append(trendDown, down[i])
		}

		if closeLine[i-1] > trendDown[i-1] {
			trend = append(trend, 1)
		} else if closeLine[i-1] < trendUp[i-1] {
			trend = append(trend, -1)
		} else {
			trend = append(trend, trend[i-1])
		}

		if trend[i] == 1 {
			tsl = append(tsl, trendUp[i])
		} else {
			tsl = append(tsl, trendDown[i])
		}
	}

	for i := 0; i < len(chartData); i++ {
		if trend[i] == 1 {
			longcondition2 = append(longcondition2, 1)
		} else {
			longcondition2 = append(longcondition2, 0)
		}

		if trend[i] == -1 {
			shortcondition2 = append(shortcondition2, 1)
		} else {
			shortcondition2 = append(shortcondition2, 0)
		}
	}

	vortexBuyP, vortexBuyM := fx.Vortex(vortexBuyPeriod, highLine, lowLine, closeLine)
	vortexSellP, vortexSellM := fx.Vortex(vortexSellPeriod, highLine, lowLine, closeLine)
	for i := 0; i < len(chartData); i++ {
		if vortexBuyP[i] > vortexBuyM[i] {
			vortexBuyPermit = append(vortexBuyPermit, 1)
		} else {
			vortexBuyPermit = append(vortexBuyPermit, 0)
		}

		if vortexSellP[i] < vortexSellM[i] {
			vortexSellPermit = append(vortexSellPermit, 1)
		} else {
			vortexSellPermit = append(vortexSellPermit, 0)
		}
	}

	for i := 0; i < len(chartData); i++ {
		if longcondition[i] == 1 && longcondition2[i] == 1 && vortexBuyPermit[i] == 1 {
			longEntry = append(longEntry, 1)
		} else {
			longEntry = append(longEntry, 0)
		}

		if shortcondition[i] == 1 && shortcondition2[i] == 1 && vortexSellPermit[i] == 1 {
			shortEntry = append(shortEntry, 1)
		} else {
			shortEntry = append(shortEntry, 0)
		}
	}

	lastSignal := 0
	for i := 0; i < len(longEntry); i++ {
		if longEntry[i] == 1 && lastSignal != 1 {
			result = append(result, 1)
			lastSignal = 1
		} else if shortEntry[i] == 1 && lastSignal != -1 {
			result = append(result, -1)
			lastSignal = -1
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
