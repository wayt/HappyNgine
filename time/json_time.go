package time

import (
	"fmt"
	"time"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {

	return []byte(fmt.Sprintf("%d", t.Unix())), nil
}

func Now() JSONTime {
	return time.Now()
}
