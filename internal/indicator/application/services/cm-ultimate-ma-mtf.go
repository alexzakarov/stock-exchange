package services

import (
	"context"
	ti "github.com/cinar/indicator"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorCMUltimateMaMTF Indicator Calculate Func
func (u *service) IndicatorCMUltimateMaMTF(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorCMUltimateMaMTF begin to work")
	}

	var closeLine []float64
	var signal int
	var result []float64
	period := 20

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	out := ti.Sma(period, closeLine)
	result = append(result, 0, 0)
	for i := 2; i < len(chartData); i++ {
		if out[i] >= out[i-2] {
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
