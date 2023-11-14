package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func (ni NullTime) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", ni.Time.Format(time.RFC3339))
	return []byte(val), nil
}

func (ni NullTime) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Time)
	ni.Valid = (err == nil)
	return err
}

func NewNullTime(t time.Time) NullTime {
	return NullTime{
		NullTime: sql.NullTime{
			Time:  t,
			Valid: true,
		},
	}
}
