package cassandra

import (
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

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

func (t JSONTime) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {

	switch info.Type() {
	case gocql.TypeTimestamp:
		if t.Unix() == 0 {
			return []byte{}, nil
		}
		x := int64(t.UTC().Unix()*1e3) + int64(t.UTC().Nanosecond()/1e6)
		return encBigInt(x), nil
	}

	return nil, errors.New(fmt.Sprintf("can not marshal JSONTime into %s", info))
}

func (t *JSONTime) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {

	switch info.Type() {
	case gocql.TypeTimestamp:
		if len(data) == 0 {
			return nil
		}
		x := decBigInt(data)
		sec := x / 1000
		nsec := (x - sec*1000) * 1000000
		*t = JSONTime{Unix(sec, nsec).In(time.UTC)}
		return nil
	}

	return errors.New(fmt.Sprintf("can not unmarshal %s", info))
}

func encBigInt(x int64) []byte {
	return []byte{byte(x >> 56), byte(x >> 48), byte(x >> 40), byte(x >> 32),
		byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)}
}

func decBigInt(data []byte) int64 {
	if len(data) != 8 {
		return 0
	}
	return int64(data[0])<<56 | int64(data[1])<<48 |
		int64(data[2])<<40 | int64(data[3])<<32 |
		int64(data[4])<<24 | int64(data[5])<<16 |
		int64(data[6])<<8 | int64(data[7])
}
