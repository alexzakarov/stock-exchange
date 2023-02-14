package services

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"main/config"
	ent "main/internal/indicator/domain/entities"
	"main/internal/indicator/domain/ports"
	"main/pkg/logger"
	"main/pkg/market_data/binance"
	"net/url"
	"reflect"
	"sync"
)

var (
	err error
)

// service Indicator service struct
type service struct {
	cfg       *config.Config
	logger    logger.Logger
	pgRepo    ports.IPostgresqlRepository
	redisRepo ports.IRedisRepository
}

// NewService Indicator domain service constructor
func NewService(cfg *config.Config, logger logger.Logger, pgRepo ports.IPostgresqlRepository, redisRepo ports.IRedisRepository) ports.IService {
	fmt.Println("test: ", &service{cfg: cfg, logger: logger, pgRepo: pgRepo, redisRepo: redisRepo})
	return &service{cfg: cfg, logger: logger, pgRepo: pgRepo, redisRepo: redisRepo}
}

func (w *service) removeDuplicateValues(Slice []ent.IndicatorIndex) map[string]int64 {
	keys := make(map[string]bool)
	list := map[string]int64{}

	for _, entry := range Slice {
		if _, value := keys[entry.AssetsSymbol]; !value {
			keys[entry.AssetsSymbol] = true
			list[entry.AssetsSymbol] = entry.AssetsId
		}
	}
	return list
}

// CalculateByInterval from database
func (w *service) CalculateByInterval(ctx context.Context, intervals string) {

	fmt.Println("CalculateByInterval begin to work")

	wgForFetchData := new(sync.WaitGroup)
	wgForCalculation := new(sync.WaitGroup)

	rows, err := w.pgRepo.ReadIndicatorIndexByInterval(ctx, intervals)
	if err != nil {
		println(err.Error())
	}
	fmt.Println(rows)
	if len(rows) > 0 {
		filteredDuplicates := w.removeDuplicateValues(rows)
		fmt.Printf("filteredDuplicates := %v \n\n", filteredDuplicates)
		w.ConsumeData(ctx, filteredDuplicates, intervals, wgForFetchData)
		wgForFetchData.Wait()
		// indicator function invoker will inject here
		for _, irow := range rows {
			wgForCalculation.Add(1)
			func(row ent.IndicatorIndex) {
				lineData, err := w.redisRepo.GetCache(ctx, row.ChannelTag)
				if err != nil {
					fmt.Printf("redis err: %s", err.Error())
				}

				values1 := w.invoke(w, row.FuncName, ctx, lineData)
				responser, err := values1[0].Interface().(ent.IndicatorCalcResponse), nil
				if err != nil {
					fmt.Printf("INVOKER ERROR : %v \n\n", err)
				}
				if w.cfg.Server.APP_DEBUG == true {
					fmt.Printf("IND: %s / PAIR: %s / INTERVAL: %s / asset: %s / assetId: %d \n ", row.FuncName, row.AssetsSymbol, row.Intervals, row.AssetsSymbol, row.AssetsId)
				}

				w.pgRepo.ChangeLockStatusByRow(ctx, row.Id)
				w.pgRepo.UpdateIndicatorResult(ctx, row, responser)
				wgForCalculation.Done()
			}(irow)
		}

		wgForCalculation.Wait()

		w.pgRepo.ReleaseAllLocks(ctx)
	}

}

func (w *service) ConsumeData(ctx context.Context, rows map[string]int64, intervals string, wgForFetchData *sync.WaitGroup) {
	if len(rows) > 0 {
		for keys, vals := range rows {
			wgForFetchData.Add(1)
			go func(key string, val int64) {

				channelTag := fmt.Sprintf("binance-%s_%s", key, intervals)
				data := &ent.IndicatorCalcRequest{Pair: key, Interval: intervals, Limit: 500}
				lineData := binance.GetKlineData(w.cfg, data.Pair, data.Interval, data.Limit)

				errOnSetCache := w.redisRepo.SetCache(ctx, channelTag, 240, lineData)
				if errOnSetCache != nil {
					println("err when write to cache ! : " + errOnSetCache.Error())
				}

				boosterRes, _ := w.BoosterStatistics(ctx, lineData)
				err := w.pgRepo.SaveBoosterStatistics(ctx, intervals, val, key, boosterRes)
				if err != nil {
					println("err when save booster statistics ! : " + err.Error())
				}
				wgForFetchData.Done()
			}(keys, vals)
		}
	}
}

func (w *service) invoke(obj any, name string, args ...any) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	return reflect.ValueOf(obj).MethodByName(name).Call(inputs)
}

func (w *service) WebsocketClient() {
	var addr = flag.String("addr", "stream.binance.com:9443", "http service address")
	u := url.URL{Scheme: "wss", Host: *addr, Path: "/btcusdt@kline_1h"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
}
