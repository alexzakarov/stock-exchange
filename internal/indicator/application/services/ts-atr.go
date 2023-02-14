package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorTsAtr Indicator Calculate Func
func (u *service) IndicatorTsAtr(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorTsAtr begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var trend []bool
	var stop []float64
	var min float64
	var max float64
	var signal int
	var result []float64

	stop = append(stop, 0)
	trend = append(trend, false)

	period := 20
	factor := 2.0
	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
	}

	for i := 1; i < len(chartData); i++ {
		currentLows := lowLine[0 : i+1]
		currentHighs := highLine[0 : i+1]
		currentCloses := closeLine[0 : i+1]
		if max == 0.0 {
			max = 0.0
		}

		if min == 0.0 {
			min = currentCloses[len(currentCloses)-1]
		}

		tr := fx.Tr(currentHighs, currentLows, currentCloses)
		if len(tr) == 0 {
			tr = append(tr, 0)
		}

		atr := fx.Rma(period, tr)
		atrlast := atr[len(atr)-1]
		atrM := atrlast * factor
		max = math.Max(max, currentCloses[len(currentCloses)-1])
		min = math.Min(min, currentCloses[len(currentCloses)-1])

		if trend[len(trend)-1] == true {
			stop = append(stop, math.Max(stop[len(stop)-1], max-atrM))
		} else {
			stop = append(stop, math.Min(stop[len(stop)-1], min+atrM))
		}

		if currentCloses[len(currentCloses)-1]-stop[len(stop)-1] > 0.0 {
			trend = append(trend, true)
		} else {
			trend = append(trend, false)
		}

		if trend[len(trend)-1] != trend[len(trend)-2] {
			max = currentCloses[len(currentCloses)-1]
			min = currentCloses[len(currentCloses)-1]
			if trend[len(trend)-1] {
				stop[len(stop)-1] = max - atrM
			} else {
				stop[len(stop)-1] = min + atrM
			}
		}

		if trend[i] {
			result = append(result, 1)
		} else {
			result = append(result, -1)
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
