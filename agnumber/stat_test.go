package agnumber

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAverage(t *testing.T) {
	testcases := []struct {
		input    []float64
		expected float64
	}{
		{
			input:    nil,
			expected: 0,
		},
		{
			input:    []float64{},
			expected: 0,
		},
		{
			input:    []float64{1},
			expected: 1,
		},
		{
			input:    []float64{1, 2},
			expected: 1.5,
		},
		{
			input:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: 5.5,
		},
	}

	for _, testcase := range testcases {
		got := Average(testcase.input)
		require.Equal(t, testcase.expected, got, "input: %v", testcase.input)
	}
}

func TestMedian(t *testing.T) {
	testcases := []struct {
		input    []float64
		expected float64
	}{
		{
			input:    nil,
			expected: 0,
		},
		{
			input:    []float64{},
			expected: 0,
		},
		{
			input:    []float64{1},
			expected: 1,
		},
		{
			input:    []float64{1, 2},
			expected: 1.5,
		},
		{
			input:    []float64{1, 2, 3},
			expected: 2,
		},
		{
			input:    []float64{1, 3, 3},
			expected: 3,
		},
		{
			input:    []float64{1, 1, 3},
			expected: 1,
		},
		{
			input:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: 5.5,
		},
	}

	for _, testcase := range testcases {
		got := Median(testcase.input)
		require.Equal(t, testcase.expected, got, "input: %v", testcase.input)
	}
}
