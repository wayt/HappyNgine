package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"strings"
)

var Sessions map[string]*gocql.Session

const (
	MAIN_KEYSPACE_ALIAS = "__main__"
)

func init() {

	Sessions = make(map[string]*gocql.Session)

	if _, err := NewKeyspaceSession(env.Get("HAPPY_CASSANDRA_KEYSPACE"), MAIN_KEYSPACE_ALIAS); err != nil {
		log.Criticalln(err)
	}
}

func NewKeyspaceSession(keyspace, alias string) (*gocql.Session, error) {

	hostsString := env.Get("CASSANDRA_PORT_9042_TCP_ADDR")
	hosts := strings.Split(hostsString, ",")

	cfg := gocql.NewCluster(hosts...)
	cfg.Keyspace = keyspace
	cfg.Port = env.GetInt("CASSANDRA_PORT_9042_TCP_PORT")

	var err error
	Sessions[alias], err = cfg.CreateSession()

	return Sessions[alias], err
}

func TimeUUID() gocql.UUID {
	return gocql.TimeUUID()
}

func ParseUUID(input string) (gocql.UUID, error) {
	return gocql.ParseUUID(input)
}

func Query(stmt string, values ...interface{}) *gocql.Query {
	return Sessions[MAIN_KEYSPACE_ALIAS].Query(stmt, values...)
}

func ExecuteBatch(batch *gocql.Batch) error {
	return Sessions[MAIN_KEYSPACE_ALIAS].ExecuteBatch(batch)
}
