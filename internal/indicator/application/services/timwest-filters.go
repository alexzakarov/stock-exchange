package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorTimWestFilters Indicator Calculate Func
func (u *service) IndicatorTimWestFilters(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorTimWestFilters begin to work")
	}

	var closeLine []float64
	var condPru []float64
	var condPrd []float64
	var signal int
	var result []float64
	period := 63

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	for i := period; i < len(chartData); i++ {
		lastCLoseInIteration := closeLine[i]
		PeriodBeforeClose := closeLine[i-period]
		if lastCLoseInIteration > PeriodBeforeClose {
			condPru = append(condPru, 1)
		} else {
			condPru = append(condPru, 0)
		}

		if lastCLoseInIteration < PeriodBeforeClose {
			condPrd = append(condPrd, 1)
		} else {
			condPrd = append(condPrd, 0)
		}
	}

	for i := 0; i < len(condPru); i++ {
		if condPru[i] == 1 {
			result = append(result, 1)
		} else if condPrd[i] == 1 {
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
