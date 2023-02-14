package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorTsAdxWs Indicator Calculate Func
func (u *service) IndicatorTsAdxWs(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorTsAdxWs begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var plus []float64
	var minus []float64
	var sum []float64
	var adx []float64
	var up []float64
	var down []float64
	var rmaUp []float64
	var rmaDown []float64
	var adxTotal []float64
	var adxAbs []float64
	var signal int
	var result []float64
	LWadxlength := 8
	LWdilength := 9

	for _, values := range chartData {
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
		closeLine = append(closeLine, values.ClosePrice)
	}

	for i := 0; i < len(chartData); i++ {
		if i == 0 {
			up = append(up, 0)
			down = append(down, 0)
		} else {
			up = append(up, chartData[i].HighPrice-chartData[i-1].HighPrice)
			down = append(down, -(chartData[i].LowPrice - chartData[i-1].LowPrice))
		}
	}

	truerange := fx.Rma(LWdilength, fx.Tr(highLine, lowLine, closeLine))
	for i := 0; i < len(up); i++ {
		if up[i] > down[i] && up[i] > 0 {
			rmaUp = append(rmaUp, up[i])
		} else {
			rmaUp = append(rmaUp, 0)
		}

		if down[i] > up[i] && down[i] > 0 {
			rmaDown = append(rmaDown, down[i])
		} else {
			rmaDown = append(rmaDown, 0)
		}
	}

	rma := fx.Rma(LWdilength, rmaUp)
	for i := 0; i < len(rma); i++ {
		plus = append(plus, fx.CheckNaN((100*rma[i])/truerange[i]))
	}

	rma2 := fx.Rma(LWdilength, rmaDown)
	for i := 0; i < len(rma2); i++ {
		minus = append(minus, fx.CheckNaN((100*rma2[i])/truerange[i]))
	}

	for i := 0; i < len(plus); i++ {
		sum = append(sum, plus[i]+minus[i])
		if sum[i] == 0 {
			adxTotal = append(adxTotal, 1)
		} else {
			adxTotal = append(adxTotal, sum[i])
		}
		adxAbs = append(adxAbs, math.Abs(plus[i]-minus[i])/adxTotal[i])
	}

	rma3 := fx.Rma(LWadxlength, adxAbs)
	for i := 0; i < len(rma3); i++ {
		adx = append(adx, rma3[i]*100)
	}

	for i := 0; i < len(adx); i++ {
		if plus[i] > minus[i] {
			result = append(result, 1)
		} else if plus[i] <= minus[i] {
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
