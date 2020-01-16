package tests

import (
	"../binance"
	"../exmo"
	"../header"
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

	dur, psec := QueriesCount(func() {
		var recency int64 = 60
		header.DefaultGetRates(pairs, markets, recency)
	}, 400, 20)
	fmt.Println(dur, psec)

}
