package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorBBStop Indicator Calculate Func
func (u *service) IndicatorBBStop(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorBBStop begin to work")
	}

	var ohlc4 []float64
	var dev []float64
	var upper []float64
	var lower []float64
	var up []float64
	var down []float64
	var signal int
	var result []float64
	period := 20
	mult := 1

	for _, values := range chartData {
		ohlc4 = append(ohlc4, (values.HighPrice+values.LowPrice+values.OpenPrice+values.ClosePrice)/4)
	}

	wma := fx.WMA(period, ohlc4)
	std := fx.StDev(period, ohlc4)
	for i := 0; i < len(std); i++ {
		dev = append(dev, std[i]*float64(mult))
	}

	for i := 0; i < len(wma); i++ {
		upper = append(upper, wma[i]+dev[i])
		lower = append(lower, wma[i]-dev[i])
	}

	crossUp := fx.CrossUp(ohlc4, upper)
	crossDown := fx.CrossDown(ohlc4, lower)
	up = append(up, 0)
	down = append(down, 0)
	for i := 1; i < len(crossUp); i++ {
		if crossUp[i] == 1 {
			up = append(up, 1)
			down = append(down, 0)
		} else if crossDown[i] == 1 {
			up = append(up, 0)
			down = append(down, 1)
		} else {
			up = append(up, up[i-1])
			down = append(down, down[i-1])
		}
	}

	for i := 1; i < len(up); i++ {
		if down[i-1] == 1 && up[i] == 1 {
			result = append(result, 1)
		} else if up[i-1] == 1 && down[i] == 1 {
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
