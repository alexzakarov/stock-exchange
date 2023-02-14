package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorRepanocha Indicator Calculate Func
func (u *service) IndicatorRepanocha(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorRepanocha begin to work")
	}

	var hlc3 []float64
	var hl2 []float64
	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var t3_rising []bool
	var t3_falling []bool
	var t3_rising2 []bool
	var t3_falling2 []bool
	var str []float64
	var sdmp []float64
	var sdmm []float64
	var pmi []float64
	var mmi []float64
	var pdi []float64
	var mdi []float64
	var dx []float64
	var signal int
	var result []float64

	for _, values := range chartData {
		hlc3 = append(hlc3, (values.HighPrice+values.LowPrice+values.ClosePrice)/3)
		hl2 = append(hl2, (values.LowPrice+values.HighPrice)/2)
		closeLine = append(closeLine, values.ClosePrice)
	}

	src_b1 := hlc3[:len(closeLine)-1]
	src_b1 = append([]float64{1}, src_b1...)
	src2_b1 := hl2[:len(closeLine)-1]
	src2_b1 = append([]float64{1}, src2_b1...)
	_t3_len := 4
	_a1 := 0.1

	src_t3 := fx.T3(hlc3, _t3_len, _a1)
	src_b1_t3 := fx.T3(src_b1, _t3_len, _a1)
	src2_t3 := fx.T3(hl2, _t3_len, _a1)
	src2_b1_t3 := fx.T3(src2_b1, _t3_len, _a1)
	for i := 0; i < len(src_t3); i++ {
		if src_t3[i] > src_b1_t3[i] {
			t3_rising = append(t3_rising, true)
		} else {
			t3_rising = append(t3_rising, false)
		}
		if src_t3[i] < src_b1_t3[i] {
			t3_falling = append(t3_falling, true)
		} else {
			t3_falling = append(t3_falling, false)
		}
		if src2_t3[i] > src2_b1_t3[i] {
			t3_rising2 = append(t3_rising2, true)
		} else {
			t3_rising2 = append(t3_rising2, false)
		}
		if src2_t3[i] < src2_b1_t3[i] {
			t3_falling2 = append(t3_falling2, true)
		} else {
			t3_falling2 = append(t3_falling2, false)
		}
	}

	adx_len := 11
	tr := fx.Tr(highLine, lowLine, closeLine)
	for i := 1; i < len(chartData); i++ {
		highLine1 := highLine[i-1]
		lowLine1 := lowLine[i-1]

		if highLine[i]-highLine1 > lowLine1-lowLine[i] {
			pmi = append(pmi, math.Max(0, highLine[i]-highLine1))
		} else {
			pmi = append(pmi, 0)
		}

		if lowLine1-lowLine[i] > highLine[i]-highLine1 {
			mmi = append(mmi, math.Max(0, lowLine1-lowLine[i]))
		} else {
			mmi = append(mmi, 0)
		}

		str = append(str, str[i-1]-str[i-1]/float64(adx_len)-str[i-1]/float64(adx_len)+tr[i])
		sdmp = append(sdmp, sdmp[i-1]-sdmp[i-1]/float64(adx_len)+pmi[i])
		sdmm = append(sdmm, sdmm[i-1]-sdmm[i-1]/float64(adx_len)+mmi[i])
		pdi = append(pdi, 100*sdmp[i]/str[i])
		mdi = append(mdi, 100*sdmm[i]/str[i])
		dx = append(dx, 100*(math.Abs(pdi[i]-mdi[i])/(pdi[i]+mdi[i])))
	}

	signal = int(fx.GetByIndexN(result, 2))

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
