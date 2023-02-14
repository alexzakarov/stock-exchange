package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"

	ti "github.com/cinar/indicator"
)

// IndicatorSqzMomLB Indicator Calculate Func
func (u *service) IndicatorSqzMomLB(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorSqzMomentum begin to work")
	}

	var highLine []float64
	var highestLine []float64
	var lowLine []float64
	var lowestLine []float64
	var closeLine []float64
	var closestLine []float64
	var avg []float64
	var avg2 []float64
	var indices []float64
	var signal int
	var result []float64
	period := 20

	for i, values := range chartData {
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
		closeLine = append(closeLine, values.ClosePrice)
		indices = append(indices, float64(i))
	}

	highestLine = ti.Max(period, highLine)
	lowestLine = ti.Min(period, lowLine)
	closestLine = fx.Sma(period, closeLine)
	for i := 0; i < len(highestLine); i++ {
		avg = append(avg, (highestLine[i]+lowestLine[i])/2)
	}

	for i := 0; i < len(avg); i++ {
		avg2 = append(avg2, closeLine[i]-((avg[i]+closestLine[i])/2))
	}

	linreg := fx.Linreg(period, 0, avg2)
	for i := 0; i < len(linreg); i++ {
		if linreg[i] > 0 {
			result = append(result, 1)
		} else if linreg[i] < 0 {
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
