package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorCoralTrendLazyBear Indicator Calculate Func
func (u *service) IndicatorCoralTrendLazyBear(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorCoralTrendLazyBear begin to work")
	}

	var closeLine []float64
	var i1 []float64
	var i2 []float64
	var i3 []float64
	var i4 []float64
	var i5 []float64
	var i6 []float64
	var bfr []float64
	var signal int
	var result []float64
	sm := 21.0
	cd := 0.4
	di := (sm - 1.0) / (2.0 + 1.0)
	c1 := 2.0 / (di + 1.0)
	c2 := 1.0 - c1
	c3 := 3.0 * (cd*cd + cd*cd*cd)
	c4 := -3.0 * (2.0*cd*cd + cd + cd*cd*cd)
	c5 := 3.0*cd + 1.0 + cd*cd*cd + 3.0*cd*cd

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	i1 = append(i1, 0)
	i2 = append(i2, 0)
	i3 = append(i3, 0)
	i4 = append(i4, 0)
	i5 = append(i5, 0)
	i6 = append(i6, 0)
	for i := 1; i < len(chartData); i++ {
		i1 = append(i1, (c1*closeLine[i] + c2*(i1[i-1])))
		i2 = append(i2, (c1*i1[i] + c2*(i2[i-1])))
		i3 = append(i3, (c1*i2[i] + c2*(i3[i-1])))
		i4 = append(i4, (c1*i3[i] + c2*(i4[i-1])))
		i5 = append(i5, (c1*i4[i] + c2*(i5[i-1])))
		i6 = append(i6, (c1*i5[i] + c2*(i6[i-1])))
	}

	for i := 0; i < len(chartData); i++ {
		bfr = append(bfr, -cd*cd*cd*i6[i]+c3*i5[i]+c4*i4[i]+c5*i3[i])
	}

	result = append(result, 0)
	for i := 1; i < len(chartData); i++ {
		if bfr[i] > bfr[i-1] {
			result = append(result, 1)
		} else if bfr[i] < bfr[i-1] {
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
