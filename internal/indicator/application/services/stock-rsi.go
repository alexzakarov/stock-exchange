package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorStockRsi Indicator Calculate Func
func (u *service) IndicatorStockRsi(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorStockRsi begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var signal int
	var result []float64
	smoothKPeriod := 3
	smoothDPeriod := 3
	periodRsi := 14
	PeriodStoch := 14

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
	}

	_, rsi := fx.RsiPeriod(periodRsi, closeLine)
	stochResult := fx.Stoch(PeriodStoch, rsi, rsi, rsi)
	smoothK := fx.Sma(smoothKPeriod, stochResult)
	smoothD := fx.Sma(smoothDPeriod, smoothK)
	crossOver := fx.CrossUp(smoothK, smoothD)
	crossUnder := fx.CrossDown(smoothK, smoothD)
	for i := 0; i < len(chartData); i++ {
		if crossOver[i] == 1 {
			result = append(result, 1)
		} else if crossUnder[i] == 1 {
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
