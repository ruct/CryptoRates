package binance

import (
	"../header"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

type Binance struct {
	cachedRates sync.Map
	updating sync.Mutex
}

func (*Binance) GetName() string {
	return "binance"
}

func (binance *Binance) GetRate(currPair header.CurrPair, recency int64) (header.Rate, error) {
	return header.DefaultGetRate(binance, currPair, recency,
		func() (header.Rate, bool) {
			rate, ok := binance.cachedRates.Load(currPair)
			if !ok {
				return header.Rate{}, ok
			}
			return rate.(header.Rate), ok
		},
		func() error {
			return binance.renew(currPair)
		}, &binance.updating)
}

func (binance *Binance) GetTradesUrl(currPair header.CurrPair) string {
	return fmt.Sprintf("https://api.binance.com/api/v1/ticker/24hr?symbol=%v%v", currPair.First, currPair.Second)
}

func (binance *Binance) processJson(currPair header.CurrPair, jsonData map[string]interface{}) error {
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

	binance.cachedRates.Store(currPair, header.Rate{
		currPair,
		buyPrice,
		sellPrice,
		time.Now().Unix(),
	})
	return nil
}

func (binance *Binance) renew(currPair header.CurrPair) error {
	return header.DefaultRenew(binance, currPair,
		func(jsonData map[string]interface{}) error {
			return binance.processJson(currPair, jsonData)
		})
}
