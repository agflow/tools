package agnumber

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat64toa(t *testing.T) {
	require.Equal(t, "12.56789", Float64toa(12.56789))
	require.Equal(t, "12.57", Float64toa(12.56789, 2))
	require.Equal(t, "12.56789", Float64toa(12.56789, 10))
	require.Equal(t, "120", Float64toa(120.00, 10))
	require.Equal(t, "120", Float64toa(120))
}

func TestAtoi(t *testing.T) {
	n, err := Atoi("     12\n")
	require.Nil(t, err)
	require.Equal(t, 12, n)

	n, err = Atoi("12")
	require.Nil(t, err)
	require.Equal(t, 12, n)

	_, err = Atoi("12.0")
	require.NotNil(t, err)

	_, err = Atoi("12,0")
	require.NotNil(t, err)

	_, err = Atoi("12c")
	require.NotNil(t, err)
}

func TestAtoi64(t *testing.T) {
	n, err := Atoi64("     12\n")
	require.Nil(t, err)
	require.Equal(t, int64(12), n)

	n, err = Atoi64("12")
	require.Nil(t, err)
	require.Equal(t, int64(12), n)

	_, err = Atoi64("12.0")
	require.NotNil(t, err)

	_, err = Atoi64("12,0")
	require.NotNil(t, err)

	_, err = Atoi64("12c")
	require.NotNil(t, err)
}

func TestMustAtoi64(t *testing.T) {
	n := MustAtoi64("     12\n")
	require.Equal(t, int64(12), n)

	n = MustAtoi64("12")
	require.Equal(t, int64(12), n)

	require.Panicsf(t, func() { MustAtoi64("12.0") }, "did not panic for 12.0")
	require.Panicsf(t, func() { MustAtoi64("12,0") }, "did not panic for 12,0")
	require.Panicsf(t, func() { MustAtoi64("12c") }, "did not panic for 12c")
}

func TestMustAtoi(t *testing.T) {
	n := MustAtoi("     12\n")
	require.Equal(t, 12, n)

	n = MustAtoi("12")
	require.Equal(t, 12, n)

	require.Panicsf(t, func() { MustAtoi("12.0") }, "did not panic for 12.0")
	require.Panicsf(t, func() { MustAtoi("12,0") }, "did not panic for 12,0")
	require.Panicsf(t, func() { MustAtoi("12c") }, "did not panic for 12c")
}

func TestAtof64(t *testing.T) {
	n, err := Atof64("    12")
	require.Nil(t, err)
	require.Equal(t, 12.0, n)

	n, err = Atof64("12")
	require.Nil(t, err)
	require.Equal(t, 12.0, n)

	n, err = Atof64("12.0")
	require.Nil(t, err)
	require.Equal(t, 12.0, n)

	_, err = Atof64("12,0")
	require.NotNil(t, err)

	_, err = Atof64("12c")
	require.NotNil(t, err)
}

func TestMustAtof64(t *testing.T) {
	n := MustAtof64("     12\n")
	require.Equal(t, float64(12), n)

	n = MustAtof64("12")
	require.Equal(t, float64(12), n)

	n = MustAtof64("12.1")
	require.Equal(t, 12.1, n)

	require.Panicsf(t, func() { MustAtof64("12,0") }, "did not panic for 12,0")
	require.Panicsf(t, func() { MustAtof64("12c") }, "did not panic for 12c")
}

func TestMinMax(t *testing.T) {
	testCases := []struct {
		input    []int
		min, max int
	}{
		{nil, math.MinInt64, math.MaxInt64},
		{[]int{}, math.MinInt64, math.MaxInt64},
		{[]int{5, 8, 45, 1, 78}, 1, 78},
		{[]int{1, 5, 8, 1, 1}, 1, 8},
		{[]int{1, 5, 8, -1, 1}, -1, 8},
		{[]int{1, 8, 8, -1, -1}, -1, 8},
	}

	for _, testCase := range testCases {
		require.Equalf(t, testCase.min, Min(testCase.input...), "input %v", testCase.input)
		require.Equalf(t, testCase.max, Max(testCase.input...), "input %v", testCase.input)
	}
}

func TestMinMaxFloat64(t *testing.T) {
	testCases := []struct {
		input    []float64
		min, max float64
	}{
		{nil, -math.MaxFloat64, math.MaxFloat64},
		{[]float64{}, -math.MaxFloat64, math.MaxFloat64},
		{[]float64{5, 8, 45, 1, 78}, 1, 78},
		{[]float64{1, 5, 8, 1, 1}, 1, 8},
		{[]float64{1, 5, 8, -1, 1}, -1, 8},
		{[]float64{1, 8, 8, -1, -1}, -1, 8},
	}

	for _, testCase := range testCases {
		require.Equalf(t, testCase.min, MinFloat64(testCase.input...), "input %v", testCase.input)
		require.Equalf(t, testCase.max, MaxFloat64(testCase.input...), "input %v", testCase.input)
	}
}

func TestDownFrom(t *testing.T) {
	testCases := []struct {
		start, end uint
		expected   []int
	}{
		{
			start:    5,
			end:      0,
			expected: []int{5, 4, 3, 2, 1},
		},
		{
			start:    4,
			end:      1,
			expected: []int{4, 3, 2},
		},
		{
			start: 2,
			end:   5,
		},
	}

	for _, testCase := range testCases {
		result := DownFrom(testCase.start, testCase.end)
		require.Equal(t, testCase.expected, result)
	}
}

func TestUpTo(t *testing.T) {
	testCases := []struct {
		start, end uint
		expected   []int
	}{
		{
			start:    2,
			end:      6,
			expected: []int{2, 3, 4, 5},
		},
		{
			start:    0,
			end:      3,
			expected: []int{0, 1, 2},
		},
		{
			start:    11,
			end:      15,
			expected: []int{11, 12, 13, 14},
		},
		{
			start: 15,
			end:   11,
		},
	}

	for _, testCase := range testCases {
		result := UpTo(testCase.start, testCase.end)
		require.Equal(t, testCase.expected, result)
	}
}

func TestSortInt64s(t *testing.T) {
	testCases := []struct {
		orig, expected []int64
	}{
		{
			orig:     []int64{3, 4, 2, 1},
			expected: []int64{1, 2, 3, 4},
		},
	}
	for _, testCase := range testCases {
		sorted := testCase.orig
		SortInt64s(sorted)
		require.Equal(t, testCase.expected, sorted)
	}
}

func TestFloorToMostSignificantDigit(t *testing.T) {
	require.Equal(t, 700.0, FloorToMostSignificantDigit(795))
	require.Equal(t, 10.0, FloorToMostSignificantDigit(11))
	require.Equal(t, 90.0, FloorToMostSignificantDigit(99))
	require.Equal(t, 200.0, FloorToMostSignificantDigit(200))
	require.Equal(t, 8.0, FloorToMostSignificantDigit(8))
	require.Equal(t, 10.0, FloorToMostSignificantDigit(10.986))
	require.Equal(t, 50.0, FloorToMostSignificantDigit(56.76232))
}

func TestAbsInt(t *testing.T) {
	require.Equal(t, 456, AbsInt(456))
	require.Equal(t, 456, AbsInt(-456))
	require.Equal(t, 10, AbsInt(10))
	require.Equal(t, 10, AbsInt(-10))
	require.Equal(t, 0, AbsInt(0))
	require.Equal(t, 200, AbsInt(-200))
}
