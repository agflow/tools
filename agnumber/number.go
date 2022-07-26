package agnumber

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

// Float64toa converts float64 value to 10-based string.
// Function takes optional argument - precision - which is described in strconv.FormatFloat
func Float64toa(x float64, precision ...int) string {
	p := -1
	if len(precision) > 0 {
		p = precision[0]
	}
	s := strconv.FormatFloat(x, 'f', p, 64)
	if strings.Count(s, ".") == 1 {
		return strings.TrimRight(strings.TrimRight(s, "0"), ".")
	}
	return s
}

// Atoi64 converts 10-based string into int64 value.
func Atoi64(s string) (int64, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.ParseInt(s, 10, 64)
	return n, errors.Wrap(err, "can't parse int")
}

// Atof64 converts 10-based string into float64 value.
func Atof64(s string) (float64, error) {
	s = strings.TrimSpace(s)
	f, err := strconv.ParseFloat(s, 64)
	return f, errors.Wrap(err, "can't parse float")
}

// MustAtoi64 is similar to Atoi64.
// But it panics if the input can't be parse as an int64
func MustAtoi64(s string) int64 {
	i, err := Atoi64(s)
	if err != nil {
		panic(err)
	}
	return i
}

// MustAtof64 is similar to Atof64.
// But it panics if the input can't be parse as an float64
func MustAtof64(s string) float64 {
	f, err := Atof64(s)
	if err != nil {
		panic(err)
	}
	return f
}

// Atoi is similar to strconv.Atoi
// but trims spaces and adds stack to errors
func Atoi(s string) (int, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s)
	return n, errors.Wrapf(err, "can't convert %q to int", s)
}

// MustAtoi converts string to integer and panics otherwise.
func MustAtoi(s string) int {
	i, err := Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

// Min finds the minimum value in int slice
func Min(values ...int) int {
	if len(values) == 0 {
		return math.MinInt64
	}
	return int(funk.Reduce(values, math.Min, math.MaxInt64).(float64))
}

// Max finds the maximum value in int slice
func Max(values ...int) int {
	if len(values) == 0 {
		return math.MaxInt64
	}
	return int(funk.Reduce(values, math.Max, math.MinInt64).(float64))
}

// MinFloat64 finds the maximum value in float64 slice
func MinFloat64(vals ...float64) float64 {
	if len(vals) == 0 {
		return -math.MaxFloat64
	}
	result := vals[0]
	for _, v := range vals {
		result = math.Min(result, v)
	}
	return result
}

// MaxFloat64 finds the maximum value in float64 slice
func MaxFloat64(vals ...float64) float64 {
	if len(vals) == 0 {
		return math.MaxFloat64
	}
	result := vals[0]
	for _, v := range vals {
		result = math.Max(result, v)
	}
	return result
}

// DownFrom generates decreasing numbers from start to end
func DownFrom(start, end uint) []int {
	s, e := int(start), int(end)
	diff := s - e
	if diff < 0 {
		return nil
	}
	nums := make([]int, diff)
	for i := range nums {
		nums[i] = s - i
	}
	return nums
}

// UpTo generates increasing numbers from start to end
func UpTo(start, end uint) []int {
	s, e := int(start), int(end)
	diff := e - s
	if diff < 0 {
		return nil
	}
	nums := make([]int, diff)
	for i := range nums {
		nums[i] = s + i
	}
	return nums
}

// SortInt64s returns a sorted slice of int64
// follows the same signature as the sort functions in the sort
// standard library sort package
func SortInt64s(s []int64) {
	sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })
}

// FloorToMostSignificantDigit rounds down to most significant digit
// e.g. 795 -> 700, 11 -> 10, 99 -> 90, 200 -> 200
func FloorToMostSignificantDigit(f float64) float64 {
	n := int(f)
	i := 0
	for n > 10 {
		n /= 10
		i++
	}
	return float64(n) * math.Pow10(i)
}

// AbsInt returns the absolute value of n.
// e.g. 795 -> 795, -11 -> 11
func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
