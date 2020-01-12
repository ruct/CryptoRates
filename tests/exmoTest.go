package tests

import (
	"../exmo"
	"../header"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var currencies = []string{
	"EXM", "USD", "EUR", "RUB", "PLN", "TRY", "UAH", "BTC", "LTC", "DOGE", "DASH", "ETH", "WAVES", "ZEC", "USDT",
	"XMR", "XRP", "KICK", "ETC", "BCH", "BTG", "EOS", "BTCZ", "DXT", "XLM", "MNX", "OMG", "TRX", "ADA", "INK", "NEO", "GAS",
	"ZRX", "GNT", "GUSD", "LSK", "XEM", "SMART", "QTUM", "HB", "DAI", "MKR", "MNC", "PTI", "ATMCASH", "ETZ", "USDC", "ROOBEE",
	"DCR", "XTZ", "ZAG", "BTT", "VLX", "HP",
}
var currPairs = []header.CurrPair{
	{"BTC", "USD"},
	{"BTC", "USDT"},
	{"ADA", "ETH"},
	{"ADA", "BTC"},
	{"DCR", "BTC"},
	{"GAS", "USD"},
	{"ETZ", "ETH"},
}

func getRandomRate(exmo *exmo.Exmo, recency int64) error {
	return getCertainRate(exmo, currPairs[rand.Intn(len(currPairs))], recency)
}

func getCertainRate(exmo *exmo.Exmo, currPair header.CurrPair, recency int64) error {
	rate, err := exmo.GetRate(currPair, recency)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(rate)
	return nil
}

func endlessGetCertainRate(exmo *exmo.Exmo, currPair header.CurrPair, recency int64) error {
	for {
		fmt.Println("cert")
		err := getCertainRate(exmo, currPair, recency)
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(time.Second*time.Duration(3+rand.Intn(5)))
	}
}
func endlessGetRandomRate(exmo *exmo.Exmo, recency int64) error {
	for {
		fmt.Println("rnd")
		err := getRandomRate(exmo, recency)
		if err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(time.Second*time.Duration(3+rand.Intn(5)))
	}
}

func ExmoTest() {
	header.Init()
	var exmo exmo.Exmo
	for i := 0; i < len(currPairs)*4; i++ {
		go endlessGetCertainRate(&exmo, currPairs[i%len(currPairs)], 40)
	}
	for i := 0; i < len(currPairs)*4; i++ {
		go endlessGetRandomRate(&exmo, 40)
	}
}