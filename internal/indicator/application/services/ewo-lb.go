package services

import (
	"context"
	ti "github.com/cinar/indicator"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorEWOLB Indicator Calculate Func
func (u *service) IndicatorEWOLB(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorEWOLB begin to work")
	}

	var calculatedData []float64
	var closeLine []float64
	var signal int
	var result []float64
	period1 := 5
	period2 := 35

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	ema5 := ti.Ema(period1, closeLine)
	ema35 := ti.Ema(period2, closeLine)

	for i := 0; i < len(chartData); i++ {
		calculatedData = append(calculatedData, ema5[i]-ema35[i])
	}

	for i := 0; i < len(calculatedData); i++ {
		if calculatedData[i] <= 0 {
			result = append(result, -1)
		} else {
			result = append(result, 1)
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
