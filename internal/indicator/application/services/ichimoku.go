package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorIchimoku Indicator Calculate Func
func (u *service) IndicatorIchimoku(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorIchimoku begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var senkouA []float64
	var ss_high []float64
	var ss_low []float64
	var tk_cross_bull []float64
	var tk_cross_bear []float64
	var cs_cross_bull []float64
	var cs_cross_bear []float64
	var price_above_kumo []float64
	var price_below_kumo []float64
	var bullish []float64
	var bearish []float64
	var signal int
	var result []float64
	ts_bars := 9
	ks_bars := 26
	ssb_bars := 52
	ss_offset := 26
	cs_offset := 26
	lastSignal := 0

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
	}

	tenkan := fx.Middle(chartData, ts_bars)
	kijun := fx.Middle(chartData, ks_bars)
	senkouB := fx.Middle(chartData, ssb_bars)
	for i := 0; i < len(chartData); i++ {
		senkouA = append(senkouA, (tenkan[i]+kijun[i])/2)
	}

	for i := ss_offset - 1; i < len(chartData); i++ {
		ss_high = append(ss_high, math.Max(senkouA[i-ss_offset+1], senkouB[i-ss_offset+1]))
		ss_low = append(ss_low, math.Min(senkouA[i-ss_offset+1], senkouB[i-ss_offset+1]))
	}

	for i := 1; i < len(tenkan); i++ {
		if tenkan[i] > kijun[i] {
			tk_cross_bull = append(tk_cross_bull, 1)
		} else {
			tk_cross_bull = append(tk_cross_bull, 0)
		}
		if tenkan[i] < kijun[i] {
			tk_cross_bear = append(tk_cross_bear, 1)
		} else {
			tk_cross_bear = append(tk_cross_bear, 0)
		}
	}

	for i := cs_offset - 1; i < len(chartData); i++ {
		if (closeLine[i] - closeLine[i-cs_offset+1]) > 0 {
			cs_cross_bull = append(cs_cross_bull, 1)
		} else {
			cs_cross_bull = append(cs_cross_bull, 0)
		}
		if (closeLine[i] - closeLine[i-cs_offset+1]) < 0 {
			cs_cross_bear = append(cs_cross_bear, 1)
		} else {
			cs_cross_bear = append(cs_cross_bear, 0)
		}
	}

	closeLine2 := closeLine[ss_offset-1:]
	tk_cross_bull2 := tk_cross_bull[cs_offset-2:]
	tk_cross_bear2 := tk_cross_bear[cs_offset-2:]
	for i := 0; i < len(closeLine2); i++ {
		if closeLine2[i] > ss_high[i] {
			price_above_kumo = append(price_above_kumo, 1)
		} else {
			price_above_kumo = append(price_above_kumo, 0)
		}
		if closeLine2[i] < ss_low[i] {
			price_below_kumo = append(price_below_kumo, 1)
		} else {
			price_below_kumo = append(price_below_kumo, 0)
		}
		if tk_cross_bull2[i] == 1 && cs_cross_bull[i] == 1 && price_above_kumo[i] == 1 {
			bullish = append(bullish, 1)
		} else {
			bullish = append(bullish, 0)
		}
		if tk_cross_bear2[i] == 1 && cs_cross_bear[i] == 1 && price_below_kumo[i] == 1 {
			bearish = append(bearish, 1)
		} else {
			bearish = append(bearish, 0)
		}

	}

	for i := 0; i < len(bullish); i++ {
		if bullish[i] == 1 && lastSignal != 1 {
			lastSignal = 1
			result = append(result, 1)
		} else if bearish[i] == 1 && lastSignal != -1 {
			lastSignal = -1
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
