package tests

import (
	"../binance"
	"../header"
)


func BinanceTest() {
	header.Init()

	var currPairs = []header.CurrPair{
		{"PERL", "USDC"},
		{"BTC", "USDT"},
		{"ADA", "ETH"},
		{"ADA", "BTC"},
		{"DCR", "BTC"},
		{"HBAR", "USDT"},
		{"XTZ", "BTC"},
	}
	var binance binance.Binance
	for i := 0; i < len(currPairs)*4; i++ {
		go endlessGetCertainRate(&binance, currPairs[i%len(currPairs)], 40)
	}
	for i := 0; i < len(currPairs)*4; i++ {
		go endlessGetRandomRate(&binance, currPairs,40)
	}
}