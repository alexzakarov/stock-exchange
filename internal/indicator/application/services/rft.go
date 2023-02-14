package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorRft Calculate Func
func (u *service) IndicatorRft(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorRft begin to work")
	}

	var closeLine []float64
	var upWard []float64
	var downWard []float64
	var longCond []float64
	var shortCond []float64
	var condIni []float64
	var signal int
	var result []float64
	period := 14.0
	mult := 2.618

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	smoothrng := fx.Smoothrng(closeLine, period, mult)
	rngfilter := fx.Rngfilt(closeLine, smoothrng)
	upWard = append(upWard, 0)
	downWard = append(downWard, 0)
	for i := 1; i < len(chartData); i++ {
		if rngfilter[i] > rngfilter[i-1] {
			upWard = append(upWard, upWard[i-1]+1)
		} else if rngfilter[i] < rngfilter[i-1] {
			upWard = append(upWard, 0)
		} else {
			upWard = append(upWard, upWard[i-1])
		}

		if rngfilter[i] < rngfilter[i-1] {
			downWard = append(downWard, downWard[i-1]+1)
		} else if rngfilter[i] > rngfilter[i-1] {
			downWard = append(downWard, 0)
		} else {
			downWard = append(downWard, downWard[i-1])
		}
	}

	longCond = append(longCond, 0)
	shortCond = append(shortCond, 0)
	for i := 1; i < len(chartData); i++ {
		if (closeLine[i] > rngfilter[i] && closeLine[i] > closeLine[i-1] && upWard[i] > 0) ||
			(closeLine[i] > rngfilter[i] && closeLine[i] < closeLine[i-1] && upWard[i] > 0) {
			longCond = append(longCond, 1)
		} else {
			longCond = append(longCond, 0)
		}

		if (closeLine[i] < rngfilter[i] && closeLine[i] < closeLine[i-1] && downWard[i] > 0) ||
			(closeLine[i] < rngfilter[i] && closeLine[i] > closeLine[i-1] && downWard[i] > 0) {
			shortCond = append(shortCond, 1)
		} else {
			shortCond = append(shortCond, 0)
		}
	}

	condIni = append(condIni, 0)
	for i := 1; i < len(chartData); i++ {
		if longCond[i] == 1 {
			condIni = append(condIni, 1)
		} else if shortCond[i] == 1 {
			condIni = append(condIni, -1)
		} else {
			condIni = append(condIni, condIni[i-1])
		}
	}

	result = append(result, 0)
	for i := 1; i < len(chartData); i++ {
		if longCond[i] == 1 && condIni[i-1] == -1 {
			result = append(result, 1)
		} else if shortCond[i] == 1 && condIni[i-1] == 1 {
			result = append(result, -1)
		} else {
			result = append(result, result[i-1])
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
