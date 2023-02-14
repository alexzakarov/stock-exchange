package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorAnnNonRepaint Indicator Calculate Func
func (u *service) IndicatorAnnNonRepaint(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorSma200 begin to work")
	}
	//TODO : MUHAMMED ile kontrol edilecek
	var buying []float64
	var shortCondition []float64
	var longFinal []float64
	var shortFinal []float64
	var lastSignal []float64
	var signal int
	//var result []float64
	threshold := 0
	closeDiff := fx.GetDiff(chartData, u.cfg, nil, "1d")

	for _, values := range closeDiff {
		if values > float64(threshold) {
			buying = append(buying, 1)
		} else {
			buying = append(buying, 0)
		}
	}

	longCondition := buying
	for i := 0; i < len(longCondition); i++ {
		if longCondition[i] == 0 {
			shortCondition = append(shortCondition, 1)
		} else {
			shortCondition = append(shortCondition, 0)
		}
	}

	longFinal = append(longFinal, 0)
	shortFinal = append(shortFinal, 0)
	lastSignal = append(lastSignal, 0)
	for i := 1; i < len(longCondition); i++ {
		if longCondition[i] == 1 && (lastSignal[i-1] == 0 || lastSignal[i-1] == -1) {
			longFinal = append(longFinal, 1)
		} else {
			longFinal = append(longFinal, 0)
		}

		if shortCondition[i] == 1 && (lastSignal[i-1] == 0 || lastSignal[i-1] == 1) {
			shortFinal = append(shortFinal, 1)
		} else {
			shortFinal = append(shortFinal, 0)
		}

		if shortCondition[i] == 1 {
			lastSignal = append(lastSignal, -1)
		} else if longFinal[i] == 1 {
			lastSignal = append(lastSignal, 1)
		} else {
			lastSignal = append(lastSignal, lastSignal[i-1])
		}
	}

	signal = int(fx.GetByIndexN(lastSignal, 1))
	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
