package binance

import (
	"../header"
	"encoding/json"
	"errors"
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

func (binance *Binance) GetRate(currPair header.CurrPair, recency int64) (header.Rate, error) {
	binance.cachedRates.Mux.Lock()
	defer binance.cachedRates.Mux.Unlock()

	cachedRate, ok := binance.cachedRates.MuxMap[currPair]
	if !ok {
		log.Println(fmt.Sprintf("binance: no such %v", currPair))
	}
	if ok && time.Now().Unix()-cachedRate.Updated > recency {
		log.Println(fmt.Sprintf("binance: need to update %v, last update was: %v", currPair, time.Now().Unix()-cachedRate.Updated))
	}

	if !ok || time.Now().Unix()-cachedRate.Updated > recency {
		was := cachedRate.Updated

		err := binance.renew(currPair)
		if err != nil {
			log.Println(err)
			return header.Rate{}, err
		}

		cachedRate, ok = binance.cachedRates.MuxMap[currPair]
		if !ok {
			err := errors.New(fmt.Sprintf("binance: incorrect pair %v", currPair))
			log.Println(err)
			return header.Rate{}, err
		}

		became := cachedRate.Updated
		log.Println(fmt.Sprintf("binance: wanted %v, was: %v, became: %v", currPair, was, became))
	}
	return cachedRate.Rate, nil
}

func (binance *Binance) renew(currPair header.CurrPair) error {
	log.Println("binance: updating ", time.Now())
	fmt.Println("binance: updating ", time.Now())

	body, err := header.GetBody(fmt.Sprintf("https://api.binance.com/api/v1/ticker/24hr?symbol=%v%v", currPair.First, currPair.Second))
	if err != nil {
		return err
	}

	bytes := []byte(body)
	var fullData map[string]interface{}
	if err := json.Unmarshal(bytes, &fullData); err != nil {
		log.Println(err)
		return err
	}

	var buyPrice, sellPrice float64 = -1, -1
	for key, value := range fullData {
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

	if binance.cachedRates.MuxMap == nil {
		binance.cachedRates.MuxMap = make(map[header.CurrPair]header.CachedRate)
	}
	binance.cachedRates.MuxMap[currPair] = header.CachedRate{
		header.Rate{currPair, buyPrice, sellPrice},
		time.Now().Unix(),
	}
	return nil
}
