package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"

	ti "github.com/cinar/indicator"
)

// IndicatorTrendFollower Indicator Calculate Func
func (u *service) IndicatorTrendFollower(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorTrendFollower begin to work")
	}

	var closeLine []float64
	var priceRange []float64
	var chan1 []float64
	var VUD1 []float64
	var VDD1 []float64
	var VCM0 []float64
	var VAR []float64
	var WWMA []float64
	var zxEMAData []float64
	var diff []float64
	var _ret []float64
	var trend []float64
	var highLine []float64
	var lowLine []float64
	var signal int
	var result []float64
	period := 20
	maPeriod := 20
	rateInp := 1
	linPeriod := 5
	rate := float64(rateInp) / 100.0
	valpha := 2.0 / (float64(maPeriod) + 1.0)
	wwalpha := 1.0 / float64(maPeriod)
	zxLag := 0

	if float64(maPeriod)/2.0 == math.Round(float64(maPeriod)/2.0) {
		zxLag = maPeriod / 2
	} else {
		zxLag = (maPeriod - 2) / 2
	}

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
		highLine = append(highLine, values.HighPrice)
		lowLine = append(lowLine, values.LowPrice)
	}

	highest := ti.Max(300, highLine)
	lowest := ti.Min(300, lowLine)
	for i := 0; i < len(chartData); i++ {
		priceRange = append(priceRange, highest[i]-lowest[i])
		chan1 = append(chan1, priceRange[i]*float64(rate))
	}

	VUD1 = append(VUD1, 0)
	VDD1 = append(VDD1, 0)
	for i := 1; i < len(chartData); i++ {
		if closeLine[i] > closeLine[i-1] {
			VUD1 = append(VUD1, closeLine[i]-closeLine[i-1])
		} else {
			VUD1 = append(VUD1, 0)
		}
		if closeLine[i] < closeLine[i-1] {
			VDD1 = append(VDD1, closeLine[i-1]-closeLine[i])
		} else {
			VDD1 = append(VDD1, 0)
		}
	}

	VUD := ti.Sum(9, VUD1)
	VDD := ti.Sum(9, VDD1)

	VCM0 = append(VCM0, 0)
	VAR = append(VAR, closeLine[0])
	WWMA = append(WWMA, 0)
	for i := 1; i < len(chartData); i++ {
		VCM0 = append(VCM0, (VUD[i]-VDD[i])/(VUD[i]+VDD[i]))
		VAR = append(VAR, (valpha*math.Abs(VCM0[i])*closeLine[i])+((1-valpha*math.Abs(VCM0[i]))*VAR[i-1]))
		WWMA = append(WWMA, (wwalpha*closeLine[i])+((1-float64(wwalpha))*WWMA[i-1]))
	}

	for i := 0; i < zxLag; i++ {
		zxEMAData = append(zxEMAData, 0)
	}
	for i := zxLag; i < len(chartData); i++ {
		zxEMAData = append(zxEMAData, closeLine[i]+(closeLine[i]-closeLine[i-zxLag]))
	}

	masrc := fx.Ema(maPeriod, closeLine)
	linreg := fx.Linreg(linPeriod, 0, masrc)
	highest2 := fx.Max(period, linreg)
	lowest2 := fx.Min(period, linreg)

	for i := 0; i < len(chartData); i++ {
		diff = append(diff, math.Abs(highest2[i]-lowest2[i]))
		if diff[i] > chan1[i] {
			if linreg[i] > (lowest2[i] + chan1[i]) {
				trend = append(trend, 1)
			} else {
				if linreg[i] < (highest2[i] - chan1[i]) {
					trend = append(trend, -1)
				} else {
					trend = append(trend, 0)
				}
			}
		} else {
			trend = append(trend, 0)
		}
		_ret = append(_ret, fx.CheckNaN((trend[i]*diff[i])/chan1[i]))
	}

	result = append(result, 0)
	for i := 1; i < len(chartData); i++ {
		if _ret[i] > 0 {
			result = append(result, 1)
		} else if _ret[i] < 0 {
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
