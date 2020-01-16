package main

import (
	"./exmo"
	"./header"
	"./tests"
	"fmt"
	"testing"
)

var markets []header.CryptoMarket
var pairs = []header.CurrPair{
	{"BTC", "USDT"},
	//{"ADA", "ETH"},
	//{"ADA", "BTC"},
	//{"DCR", "BTC"},
	//{"XTZ", "BTC"},
}

func BenchmarkMain(b *testing.B) {
	header.Init()
	defer header.CloseDB()

	var Exmo exmo.Exmo

	var pair = header.CurrPair{"BTC", "USDT"}
	Exmo.GetRate(pair, 10)
	dur, psec := tests.QueriesCount(func() {
		var recency int64 = 60
		Exmo.GetRate(pair, recency)
	}, 400, 20)
	fmt.Println(dur, psec)
	return

	markets = []header.CryptoMarket{&exmo.Exmo{}, /* &binance.Binance{}*/}

	dur, psec = tests.QueriesCount(func() {
		var recency int64 = 60
		header.DefaultGetRates(pairs, markets, recency)
	}, 400, 20)
	fmt.Println(dur, psec)

}
