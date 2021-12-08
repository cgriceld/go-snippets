package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cgriceld/go-snippets/candlesticks/domain"
	"github.com/cgriceld/go-snippets/candlesticks/generator"
	logr "github.com/sirupsen/logrus"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

type Opened map[string]*domain.Candle

func saveToFile(wg *sync.WaitGroup, period domain.CandlePeriod, in <-chan domain.Candle) <-chan domain.Candle {
	out := make(chan domain.Candle)
	fd, err := os.Create("candles_" + string(period) + ".csv")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer func() {
			if err := fd.Close(); err != nil {
				log.Fatal(err)
			}
			close(out)
			wg.Done()
		}()

		for v := range in {
			fmt.Fprintf(fd, fmt.Sprintf("%s,%v,%.6f,%.6f,%.6f,%.6f\n", v.Ticker, v.TS.Format(time.RFC3339),
				v.Open, v.High, v.Low, v.Close))
			out <- v
		}

	}()

	return out
}

func candlesToCandles(wg *sync.WaitGroup, period domain.CandlePeriod, in <-chan domain.Candle) <-chan domain.Candle {
	var start, currPeriodTS time.Time
	var err error

	opened := make(Opened)
	out := make(chan domain.Candle)

	go func() {
		defer func() {
			for _, v := range opened {
				out <- *v
			}
			close(out)
			wg.Done()
		}()

		for cand := range in {
			if start.IsZero() {
				if start, err = domain.PeriodTS(period, cand.TS); err != nil {
					log.Fatal(err)
				}
			} else {
				if currPeriodTS, err = domain.PeriodTS(period, cand.TS); err != nil {
					log.Fatal(err)
				}
				if currPeriodTS.After(start) {
					start = currPeriodTS
					for k, v := range opened {
						out <- *v
						delete(opened, k)
					}
				}
			}

			if val, ok := opened[cand.Ticker]; !ok {
				copyCandle(opened, &cand, period, start)
			} else {
				modifyCandleOnCandle(val, &cand)
			}
		}
	}()

	return out
}

func pricesToCandles1m(wg *sync.WaitGroup, period domain.CandlePeriod, in <-chan domain.Price) <-chan domain.Candle {
	var start, currPeriodTS time.Time
	var err error

	logger := logr.New()
	opened := make(Opened)
	out := make(chan domain.Candle)

	go func() {
		defer func() {
			for _, v := range opened {
				out <- *v
			}
			close(out)
			wg.Done()
		}()

		for price := range in {
			logger.Infof("%+v", price)

			if start.IsZero() {
				if start, err = domain.PeriodTS(period, price.TS); err != nil {
					log.Fatal(err)
				}
			} else {
				if currPeriodTS, err = domain.PeriodTS(period, price.TS); err != nil {
					log.Fatal(err)
				}
				if currPeriodTS.After(start) {
					start = currPeriodTS
					for k, v := range opened {
						out <- *v
						delete(opened, k)
					}
				}
			}

			if val, ok := opened[price.Ticker]; !ok {
				openCandleOnPrice(opened, &price, period, start)
			} else {
				modifyCandleOnPrice(val, price.Value)
			}
		}
	}()

	return out
}

func main() {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	prices := pg.Prices(ctx)
	wg := sync.WaitGroup{}

	wg.Add(2)
	candles1m := pricesToCandles1m(&wg, domain.CandlePeriod1m, prices)
	closed1m := saveToFile(&wg, domain.CandlePeriod1m, candles1m)

	wg.Add(2)
	candles2m := candlesToCandles(&wg, domain.CandlePeriod2m, closed1m)
	closed2m := saveToFile(&wg, domain.CandlePeriod2m, candles2m)

	wg.Add(2)
	candles10m := candlesToCandles(&wg, domain.CandlePeriod10m, closed2m)
	closed10m := saveToFile(&wg, domain.CandlePeriod10m, candles10m)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _ = range closed10m {
		}
	}()

	<-termChan
	cancel()
	wg.Wait()
}
