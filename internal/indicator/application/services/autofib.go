package services

import (
	"context"
	ti "github.com/cinar/indicator"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
)

// IndicatorAutofib Indicator Calculate Func
func (u *service) IndicatorAutofib(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorAutofib begin to work")
	}

	var highLine []float64
	var lowLine []float64
	var result []float64
	period := 150

	if len(chartData) > 0 {
		for _, values := range chartData {
			highLine = append(highLine, values.HighPrice)
			lowLine = append(lowLine, values.LowPrice)
		}

		last := chartData[len(chartData)-1]
		lower := ti.Min(period, lowLine)
		higher := ti.Max(period, highLine)
		difference := higher[len(higher)-1] - lower[len(lower)-1]
		perClose := chartData[len(chartData)-period-1].ClosePrice

		if perClose > last.ClosePrice {
			highest := higher[len(higher)-1]
			result = []float64{
				highest,
				highest - difference*0.236,
				highest - difference*0.382,
				highest - difference*0.5,
				highest - difference*0.618,
				highest - difference*0.786,
				highest - difference,
			}
		} else {
			lowest := lower[len(lower)-1]
			result = []float64{
				lowest + difference,
				lowest + difference*0.786,
				lowest + difference*0.618,
				lowest + difference*0.5,
				lowest + difference*0.382,
				lowest + difference*0.236,
				lowest,
			}
		}
	} else {
		result = []float64{}
	}

	responser = ent.IndicatorCalcResponse{
		Result:       result,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
