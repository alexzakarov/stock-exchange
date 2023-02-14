package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorHmaTrend Indicator Calculate Func
func (u *service) IndicatorHmaTrend(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorHmaTrend begin to work")
	}

	var closeLine []float64
	var hl2 []float64
	var crossUp []float64
	var crossDown []float64
	var signal int
	var result []float64
	period := 24

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		hl2 = append(hl2, fx.Hl2(values.HighPrice, values.LowPrice))
	}

	a := fx.Hma(period, hl2)
	b := fx.Hma3(period, closeLine)
	for i := 1; i < len(chartData); i++ {
		if a[i] > b[i] && a[i-1] < b[i-1] {
			crossDown = append(crossDown, 1)
		} else {
			crossDown = append(crossDown, 0)
		}
		if a[i] < b[i] && a[i-1] > b[i-1] {
			crossUp = append(crossUp, 1)
		} else {
			crossUp = append(crossUp, 0)
		}
	}

	for i := 0; i < len(crossUp); i++ {
		if crossUp[i] == 1 {
			result = append(result, 1)
		} else if crossDown[i] == 1 {
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
