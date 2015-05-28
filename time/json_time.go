package time

import (
	"fmt"
	gotime "time"
)

type JSONTime struct {
	t gotime.Time
}

func (t JSONTime) MarshalJSON() ([]byte, error) {

	return []byte(fmt.Sprintf("%d", t.t.Unix())), nil
}

func (t JSONTime) StdTime() *gotime.Time {
	return &t.t
}

func Now() JSONTime {
	return JSONTime{t: gotime.Now()}
}

func FromStdTime(stdTime gotime.Time) JSONTime {
	return JSONTime{t: stdTime}
}
