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

func (pair CurrPair) String() string {
	return fmt.Sprintf("%v_%v", pair.First, pair.Second)
}

type CryptoExchange interface {
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

const MAXPROCS = 20
func Init() {
	runtime.GOMAXPROCS(MAXPROCS)
	logInit()
	dbInit()
}

