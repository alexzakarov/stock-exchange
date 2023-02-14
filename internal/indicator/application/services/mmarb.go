package services

import (
	"context"
	ti "github.com/cinar/indicator"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"strconv"
)

// IndicatorMMARB Indicator Calculate Func
func (u *service) IndicatorMMARB(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorMMARB begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var _Signals []float64
	var signal int
	var result []float64
	ma := map[string][]float64{}
	maResult := map[string][]float64{}

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
	}

	for i := 1; i <= 18; i++ {
		ma["ma"+strconv.Itoa(i*5)] = ti.Ema(i*5, closeLine)
	}

	for i := 1; i <= 18; i++ {
		name := "ma" + strconv.Itoa(i*5)
		change := fx.Change(ma[name])
		for j := 0; j < len(change); j++ {
			if change[j] > 0 {
				maResult[name] = append(maResult[name], 1)
			} else if change[j] < 0 {
				maResult[name] = append(maResult[name], -1)
			} else {
				maResult[name] = append(maResult[name], 0)
			}
		}
	}

	for i := 0; i < len(chartData); i++ {
		sum := 0.0
		for j := 1; j <= 18; j++ {
			name := "ma" + strconv.Itoa(j*5)
			sum += maResult[name][i]
		}
		_Signals = append(_Signals, sum)
	}

	for i := 0; i < len(chartData); i++ {
		if _Signals[i] == 18 {
			result = append(result, 1)
		} else if _Signals[i] == -18 {
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
