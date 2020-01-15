package header

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

type CurrPair struct {
	First  string
	Second string
}

func (currPair CurrPair) String() string {
	return fmt.Sprintf("%v_%v", currPair.First, currPair.Second)
}

type CryptoMarket interface {
	GetName() string
	GetRate(CurrPair, int64) (Rate, error)
	GetTradesUrl(CurrPair) string
}

func logInit() {
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(f)
	log.Println("log inited")
}

const MAXPROCS = 1000
func Init() {
	runtime.GOMAXPROCS(MAXPROCS)
	logInit()
	dbInit()
}

