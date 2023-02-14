package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorQQE Indicator Calculate Func
func (u *service) IndicatorQQE(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorQQE begin to work")
	}
	var closeLine []float64
	var atrRsi []float64
	var rsiMa1 []float64
	var dar []float64
	var longBand []float64
	var longBand1 []float64
	var shortBand []float64
	var shortBand1 []float64
	var trend []float64
	var newShortBand []float64
	var newLongBand []float64
	var fastAtrRsi []float64
	var QQExlong []float64
	var QQExshort []float64
	var qqeLong []float64
	var qqeShort []float64
	var signal int
	var result []float64
	period := 14
	sf := 5
	qqe := 4.238
	wildersPeriod := period*2 - 1

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		QQExlong = append(QQExlong, 0)
		QQExshort = append(QQExshort, 0)
	}
	_, rsi := fx.RsiPeriod(period, closeLine)

	rsiMa := fx.Ema(sf, rsi)
	rsiMa1 = append(rsiMa1, 0)
	for i := 1; i < len(chartData); i++ {
		rsiMa1 = append(rsiMa1, rsiMa[i-1])
	}

	atrRsi = append(atrRsi, 0)
	for i := 1; i < len(chartData); i++ {
		atrRsi = append(atrRsi, math.Abs(rsiMa[i-1]-rsiMa[i]))
	}

	maAtrRsi := fx.Ema(wildersPeriod, atrRsi)

	darEma := fx.Ema(wildersPeriod, maAtrRsi)

	for i := 0; i < len(chartData); i++ {
		dar = append(dar, darEma[i]*qqe)
	}

	for i := 0; i < len(chartData); i++ {
		newLongBand = append(newLongBand, rsiMa[i]-dar[i])
		newShortBand = append(newShortBand, rsiMa[i]+dar[i])
	}

	longBand = append(longBand, 0)
	shortBand = append(shortBand, 0)
	for i := 1; i < len(chartData); i++ {
		if rsiMa[i-1] > longBand[i-1] && rsiMa[i] > longBand[i-1] {
			longBand = append(longBand, math.Max(longBand[i-1], newLongBand[i]))
		} else {
			longBand = append(longBand, newLongBand[i])
		}
		if rsiMa[i-1] < shortBand[i-1] && rsiMa[i] < shortBand[i-1] {
			shortBand = append(shortBand, math.Min(shortBand[i-1], newShortBand[i]))
		} else {
			shortBand = append(shortBand, newShortBand[i])
		}
	}
	longBand1 = append(longBand1, 0)
	shortBand1 = append(shortBand1, 0)
	for i := 1; i < len(chartData); i++ {
		longBand1 = append(longBand1, longBand[i-1])
		shortBand1 = append(shortBand1, shortBand[i-1])
	}
	cross1 := fx.Cross(longBand1, rsiMa)
	cross2 := fx.Cross(rsiMa, shortBand1)

	trend = append(trend, 0)
	for i := 1; i < len(chartData); i++ {
		if cross2[i] == 1 {
			trend = append(trend, 1)
		} else if cross1[i] == 1 {
			trend = append(trend, -1)
		} else {
			trend = append(trend, trend[i-1])
		}
	}

	for i := 0; i < len(chartData); i++ {
		if trend[i] == 1 {
			fastAtrRsi = append(fastAtrRsi, longBand[i])
		} else {
			fastAtrRsi = append(fastAtrRsi, shortBand[i])
		}
	}

	for i := 1; i < len(chartData); i++ {
		if fastAtrRsi[i] < rsiMa[i] {
			QQExlong[i] = QQExlong[i-1] + 1
		} else {
			QQExlong[i] = 0
		}
		if fastAtrRsi[i] > rsiMa[i] {
			QQExshort[i] = QQExshort[i-1] + 1
		} else {
			QQExshort[i] = 0
		}
	}

	qqeLong = append(qqeLong, 0)
	qqeShort = append(qqeShort, 0)
	for i := 1; i < len(chartData); i++ {
		if QQExlong[i] == 1 {
			qqeLong = append(qqeLong, 1)
		} else {
			qqeLong = append(qqeLong, 0)
		}
		if QQExshort[i] == 1 {
			qqeShort = append(qqeShort, -1)
		} else {
			qqeShort = append(qqeShort, 0)
		}
	}

	for i := 0; i < len(chartData); i++ {
		if qqeLong[i] == 1 {
			result = append(result, 1)
		} else if qqeShort[i] == -1 {
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
	result = nil
	return responser, err
}
