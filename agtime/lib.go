package agtime

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// TruncateMonth returns time with truncated month
func TruncateMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return Day(y, m, 1)
}

// TruncateDay returns time with truncated day
func TruncateDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return Day(y, m, d)
}

// Day is a convenience function to create a day
func Day(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// ClockTime interface for Clock
type ClockTime interface {
	NowUTC() NullableTime
}

// Clock struct for ClockTime
type Clock struct {
	ClockTime
}

// NowUTC returns time.Now() in NullableTime struct
func (Clock) NowUTC() NullableTime { return NewNullTime(time.Now()) }

// MockedClock returns struct for mocking Clock
type MockedClock struct {
	mock.Mock
	ClockTime
}

// NowUTC version for MockedClock
func (m *MockedClock) NowUTC() NullableTime {
	args := m.Called()
	return args.Get(0).(NullableTime)
}
