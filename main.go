package main

import (
	"./binance"
	"./header"
	"fmt"
)

var markets []header.CryptoMarket
var marketByName = make(map[string]header.CryptoMarket)

func main() {
	var binance binance.Binance
	fmt.Println(binance.GetRate(header.CurrPair{"BTC", "USDT"}, 40))
	//tests.ExmoTest()
	//select {
	//
	//}
}
