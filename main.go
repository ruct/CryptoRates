package main

import (
	"./header"
	"./tests"
)

func main() {
	header.Init()
	go tests.BinanceTest()
	go tests.ExmoTest()
	select {

	}
}
