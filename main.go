package main

type CurrPair struct {
	first  string
	second string

}
var getRates = make(map[string]func(CurrPair) string)

func init() {

}

func GetRates(site string, currPair CurrPair) (string, error) {

}

func main() {

}
