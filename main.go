package main

import (
	"./header"
	"./tests"
)

var markets []header.CryptoMarket
var marketByName = make(map[string]header.CryptoMarket)

func main() {
	tests.ExmoTest()
	select {

	}
}
