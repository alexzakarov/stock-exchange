package services

import (
	"context"
	"fmt"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"

	ti "github.com/cinar/indicator"
)

// IndicatorTrendMagic Indicator Calculate Func
func (u *service) IndicatorTrendMagic(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorTrendMagic begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var upT []float64
	var downT []float64
	var TrendMagic []float64
	var signal int
	period := 20

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
	}

	atr := ti.Sma(period, fx.Tr2(chartData))
	cci := ti.CommunityChannelIndex(period, closeLine, closeLine, closeLine)
	for i := 0; i < len(chartData); i++ {
		upT = append(upT, lowLine[i]-atr[i])
		downT = append(downT, highLine[i]+atr[i])
	}

	TrendMagic = append(TrendMagic, upT[0])
	for i := 1; i < len(cci); i++ {
		if cci[i] >= 0 {
			if upT[i] < TrendMagic[i-1] {
				TrendMagic = append(TrendMagic, TrendMagic[i-1])
			} else {
				TrendMagic = append(TrendMagic, upT[i])
			}
		} else {
			if downT[i] > TrendMagic[i-1] {
				TrendMagic = append(TrendMagic, TrendMagic[i-1])
			} else {
				TrendMagic = append(TrendMagic, downT[i])
			}
		}
	}

	for i := 0; i < len(cci); i++ {
		fmt.Println("Trend Magic: ", cci[i])
	}

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
