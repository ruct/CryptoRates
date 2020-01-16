package tests

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"../header"
)

func getRandomRate(exchange header.CryptoExchange, pairs []header.CurrPair, recency int64) error {
	return getCertainRate(exchange, pairs[rand.Intn(len(pairs))], recency)
}

func getCertainRate(exchange header.CryptoExchange, pair header.CurrPair, recency int64) error {
	rate, err := exchange.GetRate(pair, recency)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Print(rate)
	return nil
}

func endlessGetCertainRate(exchange header.CryptoExchange, pair header.CurrPair, recency int64) error {
	for {
		err := getCertainRate(exchange, pair, recency)
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(time.Second*time.Duration(10+rand.Intn(5)))
	}
}
func endlessGetRandomRate(exchange header.CryptoExchange, pairs []header.CurrPair, recency int64) error {
	for {
		err := getRandomRate(exchange, pairs, recency)
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(time.Second*time.Duration(10+rand.Intn(5)))
	}
}

