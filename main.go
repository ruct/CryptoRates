package main

import (
	"./exmo"
	"./header"
	"fmt"
	"time"
)

var markets []header.CryptoMarket
var marketByName = make(map[string]header.CryptoMarket)

func main() {
	header.Init()

	var exmo exmo.Exmo
	for {
		fmt.Println(exmo.GetRate(header.CurrPair{"BTC", "USDT"}, 20))
		time.Sleep(time.Second)
	}
}
