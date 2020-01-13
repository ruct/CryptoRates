package binance

import (
	"../header"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Binance struct {
	cachedRates header.MuxMap
}

func (*Binance) GetName() string {
	return "binance"
}
func (binance *Binance) SetRate(rate header.Rate) {
	binance.cachedRates.MuxMap[rate.CurrPair] = rate
}

func (binance *Binance) GetRate(currPair header.CurrPair, recency int64) (header.Rate, error) {
	binance.cachedRates.Mux.Lock()
	defer binance.cachedRates.Mux.Unlock()

	return header.DefaultGetRate(binance, currPair, recency,
		func() (rate header.Rate, ok bool) {
			rate, ok = binance.cachedRates.MuxMap[currPair]
			return rate, ok
		},
		func() error {
			return binance.renew(currPair)
		})
}

func (binance *Binance) GetTradesUrl(currPair header.CurrPair) string {
	return fmt.Sprintf("https://api.binance.com/api/v1/ticker/24hr?symbol=%v%v", currPair.First, currPair.Second)
}

func (binance *Binance) processJson(currPair header.CurrPair, jsonData map[string]interface{}) error {
	if binance.cachedRates.MuxMap == nil {
		binance.cachedRates.MuxMap = make(map[header.CurrPair]header.Rate)
	}

	var err error
	var buyPrice, sellPrice float64 = -1, -1
	for key, value := range jsonData {
		if key == "bidPrice" {
			value := value.(string)
			buyPrice, err = strconv.ParseFloat(value, 64)
			if err != nil {
				log.Println(fmt.Sprintf("exmo: couldn't parse buyPrice %v", value))
			}
		} else
		if key == "askPrice" {
			value := value.(string)
			sellPrice, err = strconv.ParseFloat(value, 64)
			if err != nil {
				log.Println(fmt.Sprintf("exmo: couldn't parse sellPrice %v", value))
			}
		}
	}

	binance.cachedRates.MuxMap[currPair] = header.Rate{
		currPair,
		buyPrice,
		sellPrice,
		time.Now().Unix(),
	}
	return nil
}

func (binance *Binance) renew(currPair header.CurrPair) error {
	return header.DefaultRenew(binance, currPair,
		func(jsonData map[string]interface{}) error {
			return binance.processJson(currPair, jsonData)
		})
}
