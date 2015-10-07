package sql

import (
	"github.com/streadway/simpleuuid"
	"time"
)

func TimeUUID() string {

	id, _ := simpleuuid.NewTime(time.Now())
	return id.String()
}
