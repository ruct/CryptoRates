package main

import (
	"./header"
	"./tests"
)

func main() {
	header.Init()
	defer header.CloseDB()

	go tests.BinanceTest()
	go tests.ExmoTest()
	select {}
}
