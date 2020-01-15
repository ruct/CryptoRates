package tests

import (
	"../exmo"
	"../header"
)


func ExmoTest() {
	var currPairs = []header.CurrPair{
		{"BTC", "USD"},
		{"BTC", "USDT"},
		{"ADA", "ETH"},
		{"ADA", "BTC"},
		{"DCR", "BTC"},
		{"GAS", "USD"},
		{"ETZ", "ETH"},
	}
	var exmo exmo.Exmo
	for i := 0; i < len(currPairs)*4; i++ {
		go endlessGetCertainRate(&exmo, currPairs[i%len(currPairs)], 40)
	}
	for i := 0; i < len(currPairs)*4; i++ {
		go endlessGetRandomRate(&exmo, currPairs,40)
	}
}