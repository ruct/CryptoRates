package main

import (
	"./binance"
	"./exmo"
	"./header"
	"./utils"
	"fmt"
)

var exchanges []header.CryptoExchange
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

	exchanges = []header.CryptoExchange{&exmo.Exmo{}, &binance.Binance{}}

	dur, psec := utils.QueriesCount(func() {
		var recency int64 = 60
		utils.GetRates(pairs, exchanges, recency)
	}, 400, 20)
	fmt.Println(dur, psec)

}
