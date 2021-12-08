package main

import (
	"time"

	"github.com/cgriceld/go-snippets/candlesticks/domain"
)

func openCandleOnPrice(opened Opened, price *domain.Price, period domain.CandlePeriod, start time.Time) {
	opened[price.Ticker] = &domain.Candle{
		Ticker: price.Ticker,
		Period: period,
		Open:   price.Value,
		High:   price.Value,
		Low:    price.Value,
		Close:  price.Value,
		TS:     start,
	}
}

func modifyCandleOnPrice(candle *domain.Candle, price float64) {
	if price < candle.Low {
		candle.Low = price
	}

	if price > candle.High {
		candle.High = price
	}

	candle.Close = price
}

func copyCandle(opened Opened, candle *domain.Candle, period domain.CandlePeriod, start time.Time) {
	opened[candle.Ticker] = &domain.Candle{
		Ticker: candle.Ticker,
		Period: period,
		Open:   candle.Open,
		High:   candle.High,
		Low:    candle.Low,
		Close:  candle.Close,
		TS:     start,
	}
}

func modifyCandleOnCandle(curr *domain.Candle, candle *domain.Candle) {
	if candle.Low < curr.Low {
		curr.Low = candle.Low
	}

	if candle.High > curr.High {
		curr.High = candle.High
	}

	curr.Close = candle.Close
}
