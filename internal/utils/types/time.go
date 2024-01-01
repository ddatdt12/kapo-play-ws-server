package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type NullableTime struct {
	Time  time.Time
	Valid bool
}

func (ni NullableTime) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", ni.Time.Format(time.RFC3339))
	return []byte(val), nil
}

func (ni NullableTime) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Time)
	ni.Valid = (err == nil)
	return err
}

func (ni NullableTime) MarshalBinary() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", ni.Time.Format(time.RFC3339))
	return []byte(val), nil
}

func (ni *NullableTime) UnmarshalBinary(b []byte) error {
	err := json.Unmarshal(b, ni.Time)
	ni.Valid = (err == nil)
	return err
}

func NewNullableTime(t time.Time) NullableTime {
	return NullableTime{
		Time:  t,
		Valid: true,
	}
}
