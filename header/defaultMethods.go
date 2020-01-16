package header

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

func DefaultGetRate(market CryptoMarket, currPair CurrPair, recency int64,
    getCachedRate func() (Rate, bool), renew func() error, mux *sync.Mutex) (Rate, error) {

	cachedRate, ok := getCachedRate()
	if !ok || time.Now().Unix()-cachedRate.Updated > recency {
		was := cachedRate.Updated

		mux.Lock()
		cachedRate, ok = getCachedRate()
		if !ok {
			log.Printf("%v: no such %v", market.GetName(), currPair)
		}
		if ok && time.Now().Unix()-cachedRate.Updated > recency {
			log.Printf("%v: need to update %v, last update was: %v",
				market.GetName(), currPair, time.Now().Unix()-cachedRate.Updated)
		}
		if !ok || time.Now().Unix()-cachedRate.Updated > recency {
			err := renew()
			if err != nil {
				log.Println(err)
				return Rate{}, err
			}
			cachedRate, ok = getCachedRate()
			if !ok {
				err := errors.New(fmt.Sprintf("%v: incorrect pair %v", market.GetName(), currPair))
				log.Println(err)
				mux.Unlock()
				return Rate{}, err
			}

			became := cachedRate.Updated
			log.Printf("%v: wanted %v, was: %v, became: %v", market.GetName(), currPair, was, became)
		}
		mux.Unlock()
	}
	return cachedRate, nil
}

func DefaultRenew(market CryptoMarket, currPair CurrPair,
    processJson func(map[string]interface{}) error) error {

	log.Printf("%v: updating %v", market.GetName(), time.Now())
	fullData, err := GetJson(market.GetTradesUrl(currPair))
	if err != nil {
		return err
	}

	return processJson(fullData)
}

func DefaultGetRates(pairs []CurrPair, markets []CryptoMarket, recency int64) (string, error) {
	ratesChan := make(chan FormattedRate, MAXPROCS)
	var wg sync.WaitGroup

	for i := range pairs {
		for j := range markets {
			wg.Add(1)
			go func(i int, j int, wg *sync.WaitGroup) {
				defer wg.Done()

				rate, err := markets[j].GetRate(pairs[i], recency)
				if err != nil {
					return
				}

				var fRate FormattedRate
				fRate.FromRate(markets[j], rate)
				ratesChan <- fRate
			}(i, j, &wg)
		}
	}
	wg.Wait()
	close(ratesChan)

	var rates []FormattedRate
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
