package header

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func DefaultGetRate(market CryptoMarket, currPair CurrPair, recency int64,
    getCachedRate func() (Rate, bool), renew func() error) (Rate, error) {

	cachedRate, ok := getCachedRate()
	if !ok {
		log.Println(fmt.Sprintf("%v: no such %v", market.GetName(), currPair))
	}
	if ok && time.Now().Unix()-cachedRate.Updated > recency {
		log.Println(fmt.Sprintf("%v: need to update %v, last update was: %v", market.GetName(), currPair, time.Now().Unix()-cachedRate.Updated))
	}

	if !ok || time.Now().Unix()-cachedRate.Updated > recency {
		was := cachedRate.Updated

		err := renew()
		if err != nil {
			log.Println(err)
			return Rate{}, err
		}

		cachedRate, ok = getCachedRate()
		if !ok {
			err := errors.New(fmt.Sprintf("%v: incorrect pair %v", market.GetName(), currPair))
			log.Println(err)
			return Rate{}, err
		}

		became := cachedRate.Updated
		log.Println(fmt.Sprintf("%v: wanted %v, was: %v, became: %v", market.GetName(), currPair, was, became))
	}

	err := SaveRate(market, cachedRate)
	if err != nil {
		return Rate{}, err
	}
	return cachedRate, nil
}

func DefaultRenew(market CryptoMarket, currPair CurrPair,
    processJson func(map[string]interface{}) error) error {

	log.Println(fmt.Sprintf("%v: updating %v", market.GetName(), time.Now()))
	fullData, err := GetJson(market.GetTradesUrl(currPair))
	if err != nil {
		return err
	}
	return processJson(fullData)
}
