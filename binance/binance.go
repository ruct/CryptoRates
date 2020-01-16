package binance

import (
	"../header"
	"../utils"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

type Binance struct {
	cachedRates sync.Map
	updating    sync.Mutex
}

func (*Binance) GetName() string {
	return "binance"
}

func (binance *Binance) GetRate(pair header.CurrPair, recency int64) (header.Rate, error) {
	return utils.DefaultGetRate(binance, pair, recency,
		func() (header.Rate, bool) {
			rate, ok := binance.cachedRates.Load(pair)
			if !ok {
				return header.Rate{}, ok
			}
			return rate.(header.Rate), ok
		},
		func() error {
			return binance.renew(pair)
		}, &binance.updating)
}

func (binance *Binance) GetTradesUrl(pair header.CurrPair) string {
	return fmt.Sprintf("https://api.binance.com/api/v1/ticker/24hr?symbol=%v%v", pair.First, pair.Second)
}

func (binance *Binance) processJson(pair header.CurrPair, jsonData map[string]interface{}) error {
	var err error
	var buyPrice, sellPrice float64 = -1, -1
	for key, value := range jsonData {
		if key == "bidPrice" {
			value := value.(string)
			buyPrice, err = strconv.ParseFloat(value, 64)
			if err != nil {
				log.Printf("%v: couldn't parse buyPrice %v", binance.GetName(), value)
			}
		} else
		if key == "askPrice" {
			value := value.(string)
			sellPrice, err = strconv.ParseFloat(value, 64)
			if err != nil {
				log.Printf("%v: couldn't parse sellPrice %v", binance.GetName(), value)
			}
		}
	}

	var rate = header.Rate{
		pair,
		buyPrice,
		sellPrice,
		time.Now().Unix(),
	}

	binance.cachedRates.Store(pair, rate)
	err = header.SaveRate(binance, rate)
	if err != nil {
		return err
	}
	return nil
}

func (binance *Binance) renew(pair header.CurrPair) error {
	return utils.DefaultRenew(binance, pair,
		func(jsonData map[string]interface{}) error {
			return binance.processJson(pair, jsonData)
		})
}
