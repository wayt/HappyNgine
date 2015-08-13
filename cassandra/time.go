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
		return encBigInt(t.Unix()), nil
	}

	return nil, errors.New(fmt.Sprintf("can not marshal JSONTime into %s", info))
}

func (t *JSONTime) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {

	switch info.Type() {
	case gocql.TypeTimestamp:
		*t = Unix(bytesToInt64(data), 0)
		return nil
	}

	return errors.New(fmt.Sprintf("can not unmarshal %s", info))
}

func encBigInt(x int64) []byte {
	return []byte{byte(x >> 56), byte(x >> 48), byte(x >> 40), byte(x >> 32),
		byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)}
}

func bytesToInt64(data []byte) (ret int64) {
	for i := range data {
		ret |= int64(data[i]) << (8 * uint(len(data)-i-1))
	}
	return ret
}
