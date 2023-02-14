package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorTDS Indicator Calculate Func
func (u *service) IndicatorTDS(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorTDS begin to work")
	}

	var tdsResult []float64
	var signal int
	var result []float64

	tds := fx.TDSequential(chartData)
	for i := 0; i < len(tds); i++ {
		tdsResult = append(tdsResult, float64(-tds[i].BuySetupIndex+tds[i].SellSetupIndex))
		if tdsResult[i] > 0 {
			result = append(result, 1)
		} else if tdsResult[i] < 0 {
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
	result = nil
	return responser, err
}
