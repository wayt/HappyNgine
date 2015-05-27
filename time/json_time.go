package time

import (
	"fmt"
	gotime "time"
)

type JSONTime struct {
	gotime.Time
}

func (t JSONTime) MarshalJSON() ([]byte, error) {

	return []byte(fmt.Sprintf("%d", t.Unix())), nil
}

func Now() JSONTime {
	return JSONTime{gotime.Now()}
}
