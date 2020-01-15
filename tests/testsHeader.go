package tests

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"../header"
)

func getRandomRate(market header.CryptoMarket, currPairs []header.CurrPair, recency int64) error {
	return getCertainRate(market, currPairs[rand.Intn(len(currPairs))], recency)
}

func getCertainRate(market header.CryptoMarket, currPair header.CurrPair, recency int64) error {
	rate, err := market.GetRate(currPair, recency)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Print(rate)
	return nil
}

func endlessGetCertainRate(market header.CryptoMarket, currPair header.CurrPair, recency int64) error {
	for {
		err := getCertainRate(market, currPair, recency)
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(time.Second*time.Duration(10+rand.Intn(5)))
	}
}
func endlessGetRandomRate(market header.CryptoMarket, currPairs []header.CurrPair, recency int64) error {
	for {
		err := getRandomRate(market, currPairs, recency)
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(time.Second*time.Duration(10+rand.Intn(5)))
	}
}

