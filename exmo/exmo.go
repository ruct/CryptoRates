package exmo

import (
	"../header"
)

type Exmo struct {

}
func (*Exmo) GetName() string {
	return "exmo"
}

func (exmo *Exmo) GetRate(currPair header.CurrPair) (header.Rate, error) {
	return header.Rate{}, nil
}