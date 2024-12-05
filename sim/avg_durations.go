package sim

import (
	"math/big"
	"time"
)

func AvgDurations(durations []time.Duration) time.Duration {
	accum := new(big.Int)
	for _, d := range durations {
		accum.Add(accum, big.NewInt(int64(d)))
	}
	res := accum.Div(accum, big.NewInt(int64(len(durations))))
	return time.Duration(res.Int64())
}
