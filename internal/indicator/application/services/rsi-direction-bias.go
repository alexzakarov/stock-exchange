package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"

	ti "github.com/cinar/indicator"
)

// IndicatorRsiDirectionBias Indicator Calculate Func
func (u *service) IndicatorRsiDirectionBias(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorRsiDirectionBias begin to work")
	}

	var closeLine []float64
	var ohlc4 []float64
	var top_level []float64
	var bottom_level []float64
	var bias []float64
	var lastSignal int8
	var signal int
	period := 14

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		ohlc4 = append(ohlc4, ((values.OpenPrice + values.HighPrice + values.LowPrice + values.ClosePrice) / 4))
		top_level = append(top_level, 60)
		bottom_level = append(bottom_level, 40)
	}

	_, rsi := ti.RsiPeriod(period, ohlc4)

	crossUp := fx.CrossUp(rsi, top_level)
	crossDown := fx.CrossDown(rsi, bottom_level)
	for i := 0; i < len(crossUp); i++ {
		if crossUp[i] == 1 && lastSignal != 1 {
			lastSignal = 1
			bias = append(bias, 1)
		} else if crossDown[i] == 1 && lastSignal != -1 {
			lastSignal = -1
			bias = append(bias, -1)
		} else {
			bias = append(bias, 0)
		}
	}

	signal = int(fx.GetByIndexN(bias, 2))

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
