package agtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTruncate(t *testing.T) {
	testCases := []struct {
		input, expectedDay, expectedMonth time.Time
	}{
		{
			input:         time.Date(2022, 1, 12, 3, 4, 6, 0, time.UTC),
			expectedDay:   time.Date(2022, 1, 12, 0, 0, 0, 0, time.UTC),
			expectedMonth: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:         time.Date(2022, 2, 28, 12, 4, 6, 0, time.UTC),
			expectedDay:   time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC),
			expectedMonth: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, testCase := range testCases {
		require.Equal(t, testCase.expectedDay, TruncateDay(testCase.input))
		require.Equal(t, testCase.expectedMonth, TruncateMonth(testCase.input))
	}
}
