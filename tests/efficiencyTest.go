package tests

import (
	"sync"
	"time"
)

func QueriesCount(f func(), routines int, waves int) (time.Duration, float64) {
	beg := time.Now()

	for i := 0; i < waves; i++ {
		var wg sync.WaitGroup
		for j := 0; j < routines; j++ {
			wg.Add(1)
			go func() {
				f()
				wg.Done()
			}()
		}
		wg.Wait()
	}

	dur := time.Since(beg)
	return dur, float64(routines)*float64(waves) / dur.Seconds()
}
