package main

import (
	"./binance"
	"./exmo"
	"./header"
	"./utils"
	"fmt"
)

var markets []header.CryptoMarket
var pairs = []header.CurrPair{
	{"BTC", "USDT"},
	{"ADA", "ETH"},
	{"ADA", "BTC"},
	{"DCR", "BTC"},
	{"XTZ", "BTC"},
}

func main() {
	header.Init()
	defer header.CloseDB()

	markets = []header.CryptoMarket{&exmo.Exmo{}, &binance.Binance{}}

	dur, psec := utils.QueriesCount(func() {
		var recency int64 = 60
		utils.DefaultGetRates(pairs, markets, recency)
	}, 400, 20)
	fmt.Println(dur, psec)

}
