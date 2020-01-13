package tests

import (
	"../exmo"
	"../header"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func getRandomRate(exmo header.CryptoMarket, currPairs []header.CurrPair, recency int64) error {
	return getCertainRate(exmo, currPairs[rand.Intn(len(currPairs))], recency)
}

func getCertainRate(exmo header.CryptoMarket, currPair header.CurrPair, recency int64) error {
	rate, err := exmo.GetRate(currPair, recency)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(rate)
	return nil
}

func endlessGetCertainRate(exmo header.CryptoMarket, currPair header.CurrPair, recency int64) error {
	for {
		fmt.Println("cert")
		err := getCertainRate(exmo, currPair, recency)
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(time.Second*time.Duration(3+rand.Intn(5)))
	}
}
func endlessGetRandomRate(exmo header.CryptoMarket, currPairs []header.CurrPair, recency int64) error {
	for {
		fmt.Println("rnd")
		err := getRandomRate(exmo, currPairs, recency)
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(time.Second*time.Duration(3+rand.Intn(5)))
	}
}

func ExmoTest() {
	header.Init()

	var currPairs = []header.CurrPair{
		{"BTC", "USD"},
		{"BTC", "USDT"},
		{"ADA", "ETH"},
		{"ADA", "BTC"},
		{"DCR", "BTC"},
		{"GAS", "USD"},
		{"ETZ", "ETH"},
	}
	var exmo exmo.Exmo
	for i := 0; i < len(currPairs)*4; i++ {
		go endlessGetCertainRate(&exmo, currPairs[i%len(currPairs)], 40)
	}
	for i := 0; i < len(currPairs)*4; i++ {
		go endlessGetRandomRate(&exmo, currPairs,40)
	}
}