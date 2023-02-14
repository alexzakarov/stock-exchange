package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorUCSMMLO Indicator Calculate Func
func (u *service) IndicatorUCSMMLO(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorUCSMMLO begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var range1 []float64
	var multiplier []float64
	var midline []float64
	var color []float64
	var oscillator []float64
	var a []float64
	var b []float64
	var c []float64
	var d []float64
	var z []float64
	var y []float64
	var x []float64
	var w []float64
	var signal int
	var result []float64
	period := 100
	mult := 0.125

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
	}

	hh := fx.Max(period, highLine)
	ll := fx.Min(period, lowLine)
	for i := 0; i < len(chartData); i++ {
		range1 = append(range1, hh[i]-ll[i])
		multiplier = append(multiplier, mult*range1[i])
		midline = append(midline, ll[i]+multiplier[i]*4)
		oscillator = append(oscillator, (closeLine[i]-midline[i])/(range1[i]/2))
	}

	for i := 0; i < len(oscillator); i++ {
		if oscillator[i] > 0 && oscillator[i] < mult*2 {
			a = append(a, 1)
		} else {
			a = append(a, 0)
		}

		if oscillator[i] > 0 && oscillator[i] < mult*4 {
			b = append(b, 1)
		} else {
			b = append(b, 0)
		}

		if oscillator[i] > 0 && oscillator[i] < mult*6 {
			c = append(c, 1)
		} else {
			c = append(c, 0)
		}

		if oscillator[i] > 0 && oscillator[i] < mult*8 {
			d = append(d, 1)
		} else {
			d = append(d, 0)
		}

		if oscillator[i] < 0 && oscillator[i] > -mult*2 {
			z = append(z, 1)
		} else {
			z = append(z, 0)
		}

		if oscillator[i] < 0 && oscillator[i] > -mult*4 {
			y = append(y, 1)
		} else {
			y = append(y, 0)
		}

		if oscillator[i] < 0 && oscillator[i] > -mult*6 {
			x = append(x, 1)
		} else {
			x = append(x, 0)
		}

		if oscillator[i] < 0 && oscillator[i] > -mult*8 {
			w = append(w, 1)
		} else {
			w = append(w, 0)
		}

		//Signal power Represented by color
		if a[i] == 1 {
			color = append(color, 4)
		} else if b[i] == 1 {
			color = append(color, 3)
		} else if c[i] == 1 {
			color = append(color, 2)
		} else if d[i] == 1 {
			color = append(color, 1)
		} else if z[i] == 1 {
			color = append(color, -4)
		} else if y[i] == 1 {
			color = append(color, -3)
		} else if x[i] == 1 {
			color = append(color, -2)
		} else if w[i] == 1 {
			color = append(color, -1)
		} else {
			color = append(color, 0)
		}
	}

	for i := 0; i < len(oscillator); i++ {
		if oscillator[i] > 0 {
			result = append(result, 1)
		} else if oscillator[i] < 0 {
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

	return
}
