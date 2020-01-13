package main

import (
	"./header"
	"./tests"
)

var markets []header.CryptoMarket
var marketByName = make(map[string]header.CryptoMarket)

func main() {
	//var binance binance.Binance
	//fmt.Println(binance.GetRate(header.CurrPair{"BTC", "USDT"}, 40))
	go tests.BinanceTest()
	go tests.ExmoTest()
	select {

	}
}
