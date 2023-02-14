package services

import (
	"context"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	fx "main/pkg/utils/formulas"
	"math"
)

// IndicatorMML Indicator Calculate Func
func (u *service) IndicatorMML(ctx context.Context, chartData []binance.ChartData) (responser ent.IndicatorCalcResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("IndicatorMML begin to work")
	}

	var closeLine []float64
	var result []float64
	frame := 64.0
	mult := 1.5
	log10 := math.Log(10)
	log8 := math.Log(8)
	log2 := math.Log(2)
	lookback := math.Round(frame * mult)

	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	hh := fx.Min(int(lookback), closeLine)
	ll := fx.Max(int(lookback), closeLine)
	vLow := fx.GetByIndexN(hh, 1)
	vHigh := fx.GetByIndexN(ll, 1)
	vDist := vHigh - vLow

	tmpHgh := 0.0
	if vLow < 0 {
		tmpHgh = -vLow
	} else {
		tmpHgh = vHigh
	}

	tmpLow := 0.0
	if vLow < 0 {
		tmpLow = -vLow - vDist
	} else {
		tmpLow = vLow
	}

	shift := false
	if vLow < 0 {
		shift = true
	}

	sfVar := math.Log(tmpHgh*0.4)/log10 - math.Floor(math.Log(tmpLow*0.4)/log10)

	sr := 0.0
	if tmpHgh > 25 {
		if sfVar > 0 {
			sr = math.Exp(log10 * (math.Floor(math.Log(tmpHgh*0.4)/log10) + 1))
		} else {
			sr = math.Exp(log10 * (math.Floor(math.Log(tmpHgh*0.4) / log10)))
		}
	} else {
		sr = 100 * math.Exp(log8*math.Floor(math.Log(0.005*tmpHgh)/log8))
	}
	for _, values := range chartData {
		closeLine = append(closeLine, values.ClosePrice)
	}

	N := 0.0
	nvar1 := math.Log(sr/(tmpHgh-tmpLow)) / log8
	nvar2 := nvar1 - math.Floor(nvar1)
	if nvar1 <= 0 {
		N = 0
	} else if nvar2 == 0 {
		N = math.Floor(nvar1)
	} else {
		N = math.Floor(nvar1) + 1
	}

	SI := sr * math.Exp(-N*log8)
	MCond := 0.0
	if SI > 0 {
		MCond = (tmpHgh - tmpLow) / SI
	} else {
		MCond = 0.0000001
	}

	M := math.Floor((1.0 / log2) * math.Log(MCond))
	I := math.Round(((tmpHgh + tmpLow) * 0.5) / (SI * math.Exp((M-1)*log2)))
	bot := (I - 1) * SI * math.Exp((M-1)*log2)
	top := (I + 1) * SI * math.Exp((M-1)*log2)
	doShift := tmpHgh-top > 0.25*(top-bot) || bot-tmpLow > 0.25*(top-bot)

	er := false
	if doShift {
		er = true
	}

	mm := 0.0
	nn := 0.0
	if er == false {
		mm = M
		nn = N
	} else if er && M < 2 {
		mm = M + 1
		nn = N
	} else {
		mm = 0
		nn = N - 1
	}

	finalSI := 0.0
	if er {
		finalSI = sr * math.Exp(-nn*log8)
	} else {
		finalSI = SI
	}

	finalI := 0.0
	if er {
		finalI = math.Round(((tmpHgh + tmpLow) * 0.5) / (finalSI * math.Exp((mm-1)*log2)))
	} else {
		finalI = I
	}

	fianlBot := 0.0
	if er {
		fianlBot = (finalI - 1) * finalSI * math.Exp((mm-1)*log2)
	} else {
		fianlBot = bot
	}

	finalTop := 0.0
	if er {
		finalTop = (finalI + 1) * finalSI * math.Exp((mm-1)*log2)
	} else {
		finalTop = top
	}

	increment := (finalTop - fianlBot) / 8
	absTop := 0.0
	if shift {
		absTop = -(finalTop - 3*increment)
	} else {
		absTop = finalTop + 3*increment
	}

	plus38 := absTop
	plus28 := absTop - increment
	plus18 := absTop - 2*increment
	eightEight := absTop - 3*increment
	sevenEight := absTop - 4*increment
	sixEight := absTop - 5*increment
	fiveEight := absTop - 6*increment
	fourEight := absTop - 7*increment
	threeEight := absTop - 8*increment
	twoEight := absTop - 9*increment
	oneEight := absTop - 10*increment
	zeroEight := absTop - 11*increment
	minus18 := absTop - 12*increment
	minus28 := absTop - 13*increment
	minus38 := absTop - 14*increment

	if zeroEight > fx.GetByIndexN(closeLine, 1) {
		result = append(result, 1)
	} else if eightEight < fx.GetByIndexN(closeLine, 1) {
		result = append(result, -1)
	} else {
		result = append(result, 0)
	}

	responser = ent.IndicatorCalcResponse{
		Result: []float64{
			plus38,
			plus28,
			plus18,
			eightEight,
			sevenEight,
			sixEight,
			fiveEight,
			fourEight,
			threeEight,
			twoEight,
			oneEight,
			zeroEight,
			minus18,
			minus28,
			minus38,
		},
		CalculatedAt: chartData[len(chartData)-1].CloseTime,
	}

	return
}
