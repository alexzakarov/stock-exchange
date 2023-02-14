package services

import (
	"context"
	ti "github.com/cinar/indicator"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorVSTrend Indicator Calculate Func
func (u *service) IndicatorVSTrend(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorVSTrend begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var d1 []float64
	var vI []float64
	var signal int
	lazyBear := map[string][]float64{}
	superTrend := map[string][]float64{}
	vma := 5.0
	range1 := 3.0
	period := 24.0
	k := 1.0 / vma
	lastIndex := 2
	lastBeforeIndex := 3
	pdmsVal := 0.0
	mdmsVal := 0.0
	pdiSVal := 0.0
	mdiSVal := 0.0

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
	}

	lazyBear["close"] = append(lazyBear["close"], 1, 1, 1, 1, 1)
	lazyBear["high"] = append(lazyBear["high"], 1, 1, 1, 1, 1)
	lazyBear["low"] = append(lazyBear["low"], 1, 1, 1, 1, 1)
	lazyBear["pdms"] = append(lazyBear["pdms"], 1, 1, 1, 1, 1)
	lazyBear["mdms"] = append(lazyBear["mdms"], 1, 1, 1, 1, 1)
	lazyBear["pdm"] = append(lazyBear["pdm"], 1, 1, 1, 1, 1)
	lazyBear["mdm"] = append(lazyBear["mdm"], 1, 1, 1, 1, 1)
	lazyBear["pdi"] = append(lazyBear["pdi"], 1, 1, 1, 1, 1)
	lazyBear["mdi"] = append(lazyBear["mdi"], 1, 1, 1, 1, 1)
	lazyBear["pdis"] = append(lazyBear["pdis"], 1, 1, 1, 1, 1)
	lazyBear["mdis"] = append(lazyBear["mdis"], 1, 1, 1, 1, 1)
	lazyBear["is"] = append(lazyBear["is"], 1, 1, 1, 1, 1)
	for i := 5; i < len(chartData); i++ {
		lazyBear["close"] = append(lazyBear["close"], chartData[i-1].ClosePrice)
		lazyBear["high"] = append(lazyBear["high"], chartData[i-1].HighPrice)
		lazyBear["low"] = append(lazyBear["low"], chartData[i-1].LowPrice)
		closeLastIndex := fx.GetByIndexN(closeLine[0:i+1], lastIndex)
		closeLastBeforeIndex := fx.GetByIndexN(closeLine[0:i+1], lastBeforeIndex)
		lazyBear["pdm"] = append(lazyBear["pdm"], math.Max(0, closeLastIndex-closeLastBeforeIndex))
		lazyBear["mdm"] = append(lazyBear["mdm"], math.Max(0, closeLastBeforeIndex-closeLastIndex))
		pdmsVal = (1-k)*pdmsVal + k*fx.GetByIndexN(lazyBear["pdm"], 1)
		lazyBear["pdms"] = append(lazyBear["pdms"], pdmsVal)
		mdmsVal = (1-k)*mdmsVal + k*fx.GetByIndexN(lazyBear["mdm"], 1)
		lazyBear["mdms"] = append(lazyBear["mdms"], mdmsVal)
		s := pdmsVal + mdmsVal
		pdi := pdmsVal / s
		mdi := mdmsVal / s
		pdiSVal = (1-k)*fx.GetByIndexN(lazyBear["pdis"], 1) + k*pdi
		mdiSVal = (1-k)*fx.GetByIndexN(lazyBear["mdis"], 1) + k*mdi
		lazyBear["pdis"] = append(lazyBear["pdis"], pdiSVal)
		lazyBear["mdis"] = append(lazyBear["mdis"], mdiSVal)

		d := math.Abs(fx.GetByIndexN(lazyBear["pdis"], 1) - fx.GetByIndexN(lazyBear["mdis"], 1))
		s1 := fx.GetByIndexN(lazyBear["pdis"], 1) + fx.GetByIndexN(lazyBear["mdis"], 1)
		lazyBear["is"] = append(lazyBear["is"], (1-k)*fx.GetByIndexN(lazyBear["is"], 1)+k*d/s1)
	}

	hhv := ti.Max(int(vma), lazyBear["is"])
	llv := ti.Min(int(vma), lazyBear["is"])
	for i := 0; i < len(chartData); i++ {
		d1 = append(d1, hhv[i]-llv[i])
	}

	for i := 0; i < len(chartData); i++ {
		vI = append(vI, fx.CheckNaN((lazyBear["is"][i]-llv[i])/d1[i]))
	}

	superTrend["lows"] = append(superTrend["lows"], 0)
	superTrend["highs"] = append(superTrend["highs"], 0)
	superTrend["up"] = append(superTrend["up"], 0)
	superTrend["down"] = append(superTrend["down"], 0)
	superTrend["trendUp"] = append(superTrend["trendUp"], 0)
	superTrend["trendDown"] = append(superTrend["trendDown"], 0)
	superTrend["trend"] = append(superTrend["trend"], 0)
	for i := 1; i < len(lazyBear["high"]); i++ {
		superTrend["lows"] = append(superTrend["lows"], fx.CheckNaN((1-k*vI[i])*fx.GetByIndexN(superTrend["lows"], 1)+k*vI[i]*lazyBear["low"][i]))
		superTrend["highs"] = append(superTrend["highs"], fx.CheckNaN((1-k*vI[i])*fx.GetByIndexN(superTrend["highs"], 1)+k*vI[i]*lazyBear["high"][i]))
		_, atrVal := fx.Atr(
			int(period),
			fx.SliceByIndex(highLine, 0, len(highLine)-(i+2)),
			fx.SliceByIndex(lowLine, 0, len(lowLine)-(i+2)),
			fx.SliceByIndex(closeLine, 0, len(closeLine)-(i+2)))
		superTrend["up"] = append(superTrend["up"], fx.GetByIndexN(superTrend["lows"], 1)-(fx.GetByIndexN(atrVal, 2)*range1))
		superTrend["down"] = append(superTrend["down"], fx.GetByIndexN(superTrend["highs"], 1)+(fx.GetByIndexN(atrVal, 2)*range1))

		if lazyBear["high"][i] > fx.GetByIndexN(superTrend["trendUp"], 1) {
			superTrend["trendUp"] = append(superTrend["trendUp"], math.Max(fx.GetByIndexN(superTrend["up"], 1), fx.GetByIndexN(superTrend["trendUp"], 1)))
		} else {
			superTrend["trendUp"] = append(superTrend["trendUp"], fx.GetByIndexN(superTrend["up"], 1))
		}

		if lazyBear["low"][i] < fx.GetByIndexN(superTrend["trendDown"], 1) {
			superTrend["trendDown"] = append(superTrend["trendDown"], math.Min(fx.GetByIndexN(superTrend["down"], 1), fx.GetByIndexN(superTrend["trendDown"], 1)))
		} else {
			superTrend["trendDown"] = append(superTrend["trendDown"], fx.GetByIndexN(superTrend["down"], 1))
		}

		if lazyBear["low"][i] > fx.GetByIndexN(superTrend["trendDown"], 2) {
			superTrend["trend"] = append(superTrend["trend"], 1)
		} else {
			if lazyBear["high"][i] < fx.GetByIndexN(superTrend["trendUp"], 2) {
				superTrend["trend"] = append(superTrend["trend"], -1)
			} else {
				superTrend["trend"] = append(superTrend["trend"], fx.GetByIndexN(superTrend["trend"], 1))
			}
		}
	}

	signal = int(fx.GetByIndexN(superTrend["trend"], 1))

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
