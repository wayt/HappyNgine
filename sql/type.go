package sql

import (
	gosql "database/sql"
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
