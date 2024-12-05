package sim

import (
	"cmp"
	"math"
	"sort"
	"time"
)

func p95(durations []time.Duration) time.Duration {
	return percentile(durations, 95.0)
}

type Operatable interface {
	time.Duration | ~float64 | ~int
}

func percentile[T Operatable](input []T, percent float64) (output T) {
	if len(input) == 0 {
		return
	}
	output = percentileSorted(copySort(input), percent)
	return
}

func percentileSorted[T Operatable](sortedInput []T, percent float64) (percentile T) {
	if len(sortedInput) == 0 {
		return
	}
	index := (percent / 100.0) * float64(len(sortedInput))
	i := int(math.RoundToEven(index))
	if index == float64(int64(index)) {
		percentile = (sortedInput[i-1] + sortedInput[i]) / 2.0
	} else {
		percentile = sortedInput[i-1]
	}
	return percentile
}

// copySort copies and sorts a slice ascending.
func copySort[T cmp.Ordered](input []T) []T {
	copy := copySlice(input)
	sort.Slice(copy, func(i, j int) bool {
		return copy[i] < copy[j]
	})
	return copy
}

func copySlice[T any](input []T) []T {
	output := make([]T, len(input))
	copy(output, input)
	return output
}
