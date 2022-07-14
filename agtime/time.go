package agtime

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NullableTime same as time.Time but with an extra nullable field
type NullableTime struct {
	time.Time
	Valid bool
}

// AddDate returns similar to time.Time.AddDate, but for NullableTime
func (t *NullableTime) AddDate(years, months, days int) NullableTime {
	return NewNullTime(t.Time.AddDate(years, months, days))
}

// TruncateDay returns similar to time.Time.TruncateDay, but for NullableTime
func (t *NullableTime) TruncateDay() NullableTime {
	return NewNullTime(TruncateDay(t.Time))
}

// TruncateMonth returns similar to time.Time.TruncateMonth, but for NullableTime
func (t *NullableTime) TruncateMonth() NullableTime {
	return NewNullTime(TruncateMonth(t.Time))
}

// Equal returns similar to time.Time.Equal, but for NullableTime
func (t *NullableTime) Equal(t2 NullableTime) bool {
	return t.Valid && t2.Valid && t.Time.Equal(t2.Time)
}

// After returns similar to time.Time.After, but for NullableTime
func (t *NullableTime) After(t2 NullableTime) bool {
	return t.Valid && t2.Valid && t.Time.After(t2.Time)
}

// AfterEqual returns true if time in argument is after or equal to t
func (t *NullableTime) AfterEqual(t2 NullableTime) bool {
	return t.After(t2) || t.Equal(t2)
}

// Scan implements Scanner interface
func (t *NullableTime) Scan(value interface{}) error {
	t.Time, t.Valid = value.(time.Time)
	t.Time = t.Time.UTC()
	return nil
}

// Value implements Valuer interface
func (t NullableTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.UTC(), nil
}

// MarshalJSON marshals NullableTime
func (t NullableTime) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.Format("2006-01-02T15:04:05"))
	}
	return json.Marshal(nil)
}

// NullDate returns NullableTime with truncated day
func (t NullableTime) NullDate() NullableDate {
	return NullableDate{Time: TruncateDay(t.Time), Valid: t.Valid}
}

// NewNullTime returns time in NullableTime format
func NewNullTime(t time.Time) NullableTime {
	return NullableTime{Time: t, Valid: !t.IsZero()}
}

// NullableDate is NullableTime but with truncated day
type NullableDate NullableTime

// Scan scans NullableDate
func (t *NullableDate) Scan(value interface{}) error {
	t.Time, t.Valid = value.(time.Time)
	t.Time = t.Time.UTC()
	return nil
}

// Value implements Valuer interface
func (t NullableDate) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.UTC(), nil
}

// MarshalJSON marshals NullableDate
func (t NullableDate) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.Time.Format("2006-01-02"))
	}
	return json.Marshal(nil)
}

// UnmarshalJSON unmarshals NullableDate
func (t *NullableDate) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	var err error
	if x != nil {
		t.Time, err = time.Parse("2006-01-02", *x)
		if err != nil {
			return err
		}
		t.Valid = true
	} else {
		t.Valid = false
	}
	return nil
}

// TruncateMonth returns a nullable date with truncated month
func (t *NullableDate) TruncateMonth() NullableDate {
	return NewNullDate(TruncateMonth(t.Time))
}

// After checks if nullable date `t` is after time `a`
func (t *NullableDate) After(a time.Time) bool {
	if t.Valid {
		return t.Time.After(a)
	}
	return false
}

// Before checks if nullable date `t` is before time `a`
func (t *NullableDate) Before(a time.Time) bool {
	if t.Valid {
		return t.Time.Before(a)
	}
	return false
}

// Equal checks if nullable date `t` is equal time `a`
func (t *NullableDate) Equal(a time.Time) bool {
	if t.Valid {
		return t.Time.Equal(a)
	}
	return false
}

// AftEq checks if nullable date `t` is after or equal the time `a`
func (t *NullableDate) AftEq(a time.Time) bool {
	if t.Valid {
		return t.After(a) || t.Equal(a)
	}
	return false
}

// BfEq checks if nullable date `t` is before or equal the time `a`
func (t *NullableDate) BfEq(a time.Time) bool {
	if t.Valid {
		return t.Before(a) || t.Equal(a)
	}
	return false
}

// Between checks if nullable date `t` is between the times `a` and `b`
func (t *NullableDate) Between(a, b time.Time) bool {
	return t.AftEq(a) && t.BfEq(b)
}

// NewNullDate returns time in NullableDate format
func NewNullDate(t time.Time) NullableDate {
	return NullableDate{Time: TruncateDay(t), Valid: !t.IsZero()}
}
