package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
)

// IndicatorAlphaTrend Indicator Calculate Func
func (u *service) IndicatorAlphaTrend(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorAlphaTrend begin to work")
	}

	var closeLine []float64
	var highLine []float64
	var lowLine []float64
	var volume []float64
	var upT []float64
	var downT []float64
	var AlphaTrend []float64
	var AlphaTrend2 []float64
	var BuySignal []float64
	var SellSignal []float64
	var K1 []float64
	var K2 []float64
	var O1 []float64
	var O2 []float64
	var result []float64
	var signal int
	coeff := 1
	ap := 14

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
		volume = append(volume, values.Volume)
	}

	atr := fx.Sma(ap, fx.Tr2(chartData))
	mfi := fx.Mfi(ap, chartData)
	for i := 0; i < len(lowLine); i++ {
		upT = append(upT, lowLine[i]-atr[i]*float64(coeff))
		downT = append(downT, highLine[i]+atr[i]*float64(coeff))
	}

	AlphaTrend = append(AlphaTrend, upT[0], upT[1])
	AlphaTrend2 = append(AlphaTrend2, upT[0], upT[1])
	for i := 2; i < len(mfi); i++ {
		if mfi[i] >= 50 {
			if upT[i] < AlphaTrend[i-1] {
				AlphaTrend = append(AlphaTrend, AlphaTrend[i-1])
			} else {
				AlphaTrend = append(AlphaTrend, upT[i])
			}
		} else {
			if downT[i] > AlphaTrend[i-1] {
				AlphaTrend = append(AlphaTrend, AlphaTrend[i-1])
			} else {
				AlphaTrend = append(AlphaTrend, downT[i])
			}
		}
		AlphaTrend2 = append(AlphaTrend2, AlphaTrend[i-2])
	}

	BuySignal = append(BuySignal, 0)
	SellSignal = append(SellSignal, 0)
	for i := 1; i < len(AlphaTrend); i++ {
		if AlphaTrend[i] > AlphaTrend2[i] && AlphaTrend[i-1] <= AlphaTrend2[i-1] {
			BuySignal = append(BuySignal, 1)
		} else {
			BuySignal = append(BuySignal, 0)
		}
		if AlphaTrend[i] < AlphaTrend2[i] && AlphaTrend[i-1] >= AlphaTrend2[i-1] {
			SellSignal = append(SellSignal, 1)
		} else {
			SellSignal = append(SellSignal, 0)
		}
	}

	K1 = append(K1, 0)
	K2 = append(K2, 0)
	O1 = append(O1, 0)
	O2 = append(O2, 0)
	for i := 1; i < len(AlphaTrend); i++ {
		if BuySignal[i] == 0 {
			K1 = append(K1, K1[i-1]+1)
		} else {
			K1 = append(K1, 0)
		}
		if SellSignal[i] == 0 {
			K2 = append(K2, K2[i-1]+1)
		} else {
			K2 = append(K2, 0)
		}
		if BuySignal[i-1] == 0 {
			O1 = append(O1, O1[i-1]+1)
		} else {
			O1 = append(O1, 0)
		}
		if SellSignal[i-1] == 0 {
			O2 = append(O2, O2[i-1]+1)
		} else {
			O2 = append(O2, 0)
		}
	}

	for i := 0; i < len(AlphaTrend); i++ {
		if BuySignal[i] == 1 && O1[i] > K2[i] {
			result = append(result, 1)
		} else if SellSignal[i] == 1 && O2[i] > K1[i] {
			result = append(result, -1)
		} else {
			result = append(result, 0)
		}
	}

	signal = int(fx.GetByIndexN(result, 1))

	responser = ent.IndicatorCalcResponse{
		Signal:       signal,
		Result:       nil,
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
