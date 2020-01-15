package main

import (
	"./binance"
	"./exmo"
	"./header"
	"./tests"
	"fmt"
)


var markets []header.CryptoMarket
var pairs = []header.CurrPair{
	//{"PERL", "USDC"},
	{"BTC", "USDT"},
	{"ADA", "ETH"},
	{"ADA", "BTC"},
	{"DCR", "BTC"},
	//{"HBAR", "USDT"},
	{"XTZ", "BTC"},
}

func main() {
	header.Init()
	defer header.CloseDB()

	markets = []header.CryptoMarket{&exmo.Exmo{}, &binance.Binance{}}

	dur, psec := tests.QueriesCount(func () {
		var recency int64 = 60
		header.DefaultGetRates(pairs, markets, recency)
	}, 40, 10)
	fmt.Println(dur, psec)
}
