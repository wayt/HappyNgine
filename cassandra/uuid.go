package cassandra

import (
	"github.com/gocql/gocql"
)

func TimeUUID() gocql.UUID {
	return gocql.TimeUUID()
}

func ParseUUID(input string) (gocql.UUID, error) {
	return gocql.ParseUUID(input)
}
