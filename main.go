package main

import (
	"./binance"
	"./exmo"
	"./header"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

func getRates(pairs []header.CurrPair, markets []header.CryptoMarket, recency int64) (string, error) {
	ratesChan := make(chan header.FormattedRate, header.MAXPROCS)
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

				var fRate header.FormattedRate
				fRate.FromRate(markets[j], rate)
				ratesChan <- fRate
				fmt.Println("ended ", markets[j].GetName(), pairs[i])
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

var markets []header.CryptoMarket
var pairs = []header.CurrPair{
	{"PERL", "USDC"},
	{"BTC", "USDT"},
	{"ADA", "ETH"},
	{"ADA", "BTC"},
	{"DCR", "BTC"},
	{"HBAR", "USDT"},
	{"XTZ", "BTC"},
}

func main() {
	header.Init()
	defer header.CloseDB()

	markets = []header.CryptoMarket{&exmo.Exmo{}, &binance.Binance{}}

	var recency int64 = 40
	for {
		var s, _ = getRates(pairs, markets, recency)
		fmt.Println(s)
		time.Sleep(time.Second*time.Duration(recency)/2)
	}
}
