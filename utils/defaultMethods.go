package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
	"../header"
)

func DefaultGetRate(exchange header.CryptoExchange, pair header.CurrPair, recency int64,
    getCachedRate func() (header.Rate, bool), renew func() error, mux *sync.Mutex) (header.Rate, error) {

	cachedRate, ok := getCachedRate()
	if !ok || time.Now().Unix()-cachedRate.Updated > recency {
		was := cachedRate.Updated

		mux.Lock()
		cachedRate, ok = getCachedRate()
		if !ok {
			log.Printf("%v: no such %v", exchange.GetName(), pair)
		}
		if ok && time.Now().Unix()-cachedRate.Updated > recency {
			log.Printf("%v: need to update %v, last update was: %v",
				exchange.GetName(), pair, time.Now().Unix()-cachedRate.Updated)
		}
		if !ok || time.Now().Unix()-cachedRate.Updated > recency {
			err := renew()
			if err != nil {
				log.Println(err)
				return header.Rate{}, err
			}
			cachedRate, ok = getCachedRate()
			if !ok {
				err := errors.New(fmt.Sprintf("%v: incorrect pair %v", exchange.GetName(), pair))
				log.Println(err)
				mux.Unlock()
				return header.Rate{}, err
			}

			became := cachedRate.Updated
			log.Printf("%v: wanted %v, was: %v, became: %v", exchange.GetName(), pair, was, became)
		}
		mux.Unlock()
	}
	return cachedRate, nil
}

func DefaultRenew(exchange header.CryptoExchange, pair header.CurrPair,
    processJson func(map[string]interface{}) error) error {

	log.Printf("%v: updating %v", exchange.GetName(), time.Now())
	fullData, err := GetJson(exchange.GetTradesUrl(pair))
	if err != nil {
		return err
	}

	return processJson(fullData)
}

func GetRates(pairs []header.CurrPair, exchanges []header.CryptoExchange, recency int64) (string, error) {
	ratesChan := make(chan header.FormattedRate, header.MAXPROCS)
	var wg sync.WaitGroup

	for i := range pairs {
		for j := range exchanges {
			wg.Add(1)
			go func(i int, j int, wg *sync.WaitGroup) {
				defer wg.Done()

				rate, err := exchanges[j].GetRate(pairs[i], recency)
				if err != nil {
					return
				}

				var fRate header.FormattedRate
				fRate.FromRate(exchanges[j], rate)
				ratesChan <- fRate
			}(i, j, &wg)
		}
	}
	wg.Wait()
	close(ratesChan)

	var rates []header.FormattedRate
	for rate := range ratesChan {
		rates = append(rates, rate)
	}

	bytes, err := json.Marshal(rates)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(bytes), err
}
