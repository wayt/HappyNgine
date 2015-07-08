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
	cfg.Consistency = gocql.One
	if timeout := env.GetInt("HAPPY_CASSANDRA_TIMEOUT"); timeout > 0 {
		cfg.Timeout = time.Duration(timeout)
	}

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
	return AQuery(MAIN_KEYSPACE_ALIAS, stmt, values...)
}

func AQuery(alias, stmt string, values ...interface{}) *gocql.Query {
	return Sessions[alias].Query(stmt, values...)
}

func ExecuteBatch(batch *gocql.Batch) error {
	return AExecuteBatch(MAIN_KEYSPACE_ALIAS, batch)
}

func AExecuteBatch(alias string, batch *gocql.Batch) error {
	return Sessions[alias].ExecuteBatch(batch)
}
