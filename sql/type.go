package sql

import (
	gosql "database/sql"
	// "database/sql/driver"
	// "errors"
	"fmt"
	"time"
)

type NullString struct {
	gosql.NullString
}

func (s NullString) MarshalJSON() ([]byte, error) {
	return []byte(s.String), nil
}

func (s NullString) UnmarshalJSON(data []byte) error {
	s.String = string(data)
	s.Valid = s.String != ""
	return nil
}

func (s NullString) Len() int {
	return len(s.String)
}

type JSONTime struct {
	time.Time
}

func Now() JSONTime {
	return JSONTime{time.Now()}
}

func Unix(sec, nsec int64) JSONTime {
	return JSONTime{time.Unix(sec, nsec)}
}

func (t JSONTime) MarshalJSON() ([]byte, error) {

	return []byte(fmt.Sprintf("%d", t.Unix())), nil
}

// func (t JSONTime) Value() (driver.Value, error) {
//
// 	return t, nil
// }
//
// func (t *JSONTime) Scan(src interface{}) error {
//
// 	if src == nil {
// 		t.Time = Unix(0, 0)
// 		return nil
// 	}
//
// 	switch t := src.(type) {
// 	case time.Time:
// 		*t, ok = src.(time.Time)
// 		return errors.New("fail to cast src to time.Time")
// 	}
//
// 	return nil
// }
