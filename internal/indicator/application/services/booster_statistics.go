package services

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	ent "main/internal/indicator/domain/entities"
	"main/pkg/market_data/binance"
	"sync"
)

// Booster Statistics Calculate Func
func (u *service) BoosterStatistics(ctx context.Context, chartData []binance.ChartData) (responser ent.BoosterStatisticsResponse, err error) {
	if u.cfg.Server.APP_DEBUG == true {
		println("Booster statistics begin to work")
	}

	var calcFibo ent.IndicatorCalcResponse
	var calcMMath ent.IndicatorCalcResponse
	var calcSilver ent.IndicatorCalcResponse
	var calcGolden ent.IndicatorCalcResponse
	var calcRSI ent.IndicatorCalcResponse
	var calcMACD ent.IndicatorCalcResponse
	var calcTds ent.IndicatorCalcResponse
	var calcDeath ent.IndicatorCalcResponse

	wgForCalculation := new(sync.WaitGroup)

	wgForCalculation.Add(1)
	go func() {
		calcFibo, _ = u.IndicatorAutofib(ctx, chartData)
		wgForCalculation.Done()
	}()

	wgForCalculation.Add(1)
	go func() {
		calcMMath, _ = u.IndicatorMML(ctx, chartData)
		wgForCalculation.Done()
	}()

	wgForCalculation.Add(1)
	go func() {
		calcRSI, _ = u.IndicatorRsi(ctx, chartData)
		wgForCalculation.Done()
	}()

	wgForCalculation.Add(1)
	go func() {
		calcMACD, _ = u.IndicatorMacd(ctx, chartData)
		wgForCalculation.Done()
	}()

	wgForCalculation.Add(1)
	go func() {
		calcTds, _ = u.IndicatorTDS(ctx, chartData)
		wgForCalculation.Done()
	}()

	wgForCalculation.Add(1)
	go func() {
		calcGolden, _ = u.IndicatorGoldenCross(ctx, chartData)
		wgForCalculation.Done()
	}()

	wgForCalculation.Add(1)
	go func() {
		calcSilver, _ = u.IndicatorSilverCross(ctx, chartData)
		wgForCalculation.Done()
	}()

	wgForCalculation.Add(1)
	go func() {
		calcDeath, _ = u.IndicatorDeathCross(ctx, chartData)
		wgForCalculation.Done()
	}()

	wgForCalculation.Wait()

	responser = ent.BoosterStatisticsResponse{
		Line: ent.Line{
			Rsi:         calcRSI.Signal,
			Macd:        calcMACD.Signal,
			Tds:         calcTds.Signal,
			GoldenCross: calcGolden.Signal,
			SilverCross: calcSilver.Signal,
			DeathCross:  calcDeath.Signal,
		},
		Fibonacci: calcFibo.Result,
		MMath:     calcMMath.Result,
	}
	leo, _ := json.Marshal(responser)
	fmt.Printf("Booster statistics response: %s", leo)

	return
}
