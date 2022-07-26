package sql

import (
	q "database/sql"
	"encoding/json"
	"errors"
	"reflect"
)

// NullableString returns similar to NullString from database/sql package
type NullableString struct {
	q.NullString
}

// NewNullString creates a valid NullableString
func NewNullString(s string) NullableString {
	return NullableString{NullString: q.NullString{String: s, Valid: true}}
}

// MarshalJSON marshals NullableString
func (v NullableString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON unmarshals NullableString
func (v *NullableString) UnmarshalJSON(data []byte) error { // nolint: dupl
	// Unmarshalling into a pointer will let us detect null
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x == nil || *x == "NaN" {
		v.Valid = false
		return nil
	}

	v.Valid = true
	v.String = *x
	return nil
}

// NullableStringBind binds to nullable string
func NullableStringBind(v reflect.Value) interface{} {
	n, ok := v.Interface().(NullableString)
	if !ok || !n.Valid {
		return ""
	}
	return n.String
}

// NullableInt returns similar to NullInt64 from database/sql package
type NullableInt struct {
	q.NullInt64
}

// NewNullInt returns new NullableInt
func NewNullInt(v int64) NullableInt {
	return NullableInt{NullInt64: q.NullInt64{Int64: v, Valid: true}}
}

// MarshalJSON marshals NullableInt
func (v NullableInt) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON unmarshal NullableInt
func (v *NullableInt) UnmarshalJSON(data []byte) error { // nolint: dupl
	// Unmarshalling into a pointer will let us detect null
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

// NullableIntBind binds to nullableInt
func NullableIntBind(v reflect.Value) interface{} {
	n, ok := v.Interface().(NullableInt)
	if !ok || !n.Valid {
		return 0
	}
	return n.Int64
}

// NullableFloat returns similar to NullFloat64 from database/sql package
type NullableFloat struct {
	q.NullFloat64
}

// MarshalJSON marshals NullableFloat
func (v NullableFloat) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON unmarshal NullableFloat
func (v *NullableFloat) UnmarshalJSON(data []byte) error { // nolint: dupl
	// Unmarshalling into a pointer will let us detect null
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Float64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

// Ints is array of int64
type Ints []int64

// Scan scans Ints
func (ns *Ints) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &ns)
	case nil:
		return json.Unmarshal([]byte("null"), &ns)
	}
	return errors.New("type assertion to []byte failed")
}

// Strings is array of string
type Strings []string

// Scan scans Strings
func (ss *Strings) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &ss)
	case nil:
		return json.Unmarshal([]byte("null"), &ss)
	}
	return errors.New("type assertion []byte failed")
}

// Floats is array of float64
type Floats []float64

// Scan scans Floats
func (fs *Floats) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &fs)
}

// NewNullFloat creates a valid NullableFloat
func NewNullFloat(f float64) NullableFloat {
	return NullableFloat{NullFloat64: q.NullFloat64{Float64: f, Valid: true}}
}
