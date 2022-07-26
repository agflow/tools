package agnumber

import (
	"sort"

	"github.com/thoas/go-funk"
)

// Average finds the average of the given numbers
func Average(ls []float64) float64 {
	if len(ls) == 0 {
		return 0
	}
	return funk.SumFloat64(ls) / float64(len(ls))
}

// Median finds the median of the given numbers
func Median(ls []float64) float64 {
	if len(ls) == 0 {
		return 0
	}
	sort.Float64s(ls)
	l := len(ls)
	if l%2 == 1 {
		return ls[l/2]
	}
	return (ls[l/2] + ls[l/2-1]) / 2
}
