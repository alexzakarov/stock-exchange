package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"

	ti "github.com/cinar/indicator"
)

// IndicatorDeathCross Indicator Calculate Func
func (u *service) IndicatorDeathCross(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorDeathCross begin to work")
	}

	var closeLine []float64
	var signal int

	period1 := 50
	period3 := 200

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	sma50 := ti.Sma(period1, closeLine)
	sma200 := ti.Sma(period3, closeLine)

	cross := fx.Cross(sma50, sma200)

	signal = int(fx.GetByIndexN(cross, 2))

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
