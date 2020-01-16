package tests

import (
	"../binance"
	"../exmo"
	"../header"
	"../utils"
	"fmt"
	"testing"
)

var markets []header.CryptoMarket
var pairs = []header.CurrPair{
	{"BTC", "USDT"},
	{"ADA", "ETH"},
	{"ADA", "BTC"},
	{"DCR", "BTC"},
	{"XTZ", "BTC"},
}

func BenchmarkMain(b *testing.B) {
	header.Init()
	defer header.CloseDB()

	markets = []header.CryptoMarket{&exmo.Exmo{}, &binance.Binance{}}

	dur, psec := utils.QueriesCount(func() {
		var recency int64 = 60
		utils.DefaultGetRates(pairs, markets, recency)
	}, 400, 200)
	fmt.Println(dur, psec)

}
