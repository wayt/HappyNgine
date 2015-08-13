package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"strings"
	"time"
)

var Sessions map[string]*gocql.Session

var debugQueries bool = env.GetBool("HAPPY_CASSANDRA_DEBUG")

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
		cfg.Timeout = time.Duration(timeout) * time.Millisecond
	}

	if numConns := env.GetInt("HAPPY_CASSANDRA_NUM_CONNS"); numConns > 0 {
		cfg.NumConns = numConns
	}

	if numStreams := env.GetInt("HAPPY_CASSANDRA_NUM_STREAMS"); numStreams > 0 {
		cfg.NumStreams = numStreams
	}

	if username := env.Get("HAPPY_CASSANDRA_USERNAME"); len(username) > 0 {
		cfg.Authenticator = gocql.PasswordAuthenticator{
			Username: username,
			Password: env.Get("HAPPY_CASSANDRA_PASSWORD"),
		}
	}

	var err error
	Sessions[alias], err = cfg.CreateSession()

	return Sessions[alias], err
}

func Query(stmt string, values ...interface{}) *gocql.Query {
	return AQuery(MAIN_KEYSPACE_ALIAS, stmt, values...)
}

func AQuery(alias, stmt string, values ...interface{}) *gocql.Query {
	if debugQueries {
		log.Debugln("cassandra.AQuery[", alias, "]:", stmt)
	}
	return Sessions[alias].Query(stmt, values...)
}

func ExecuteBatch(batch *gocql.Batch) error {
	return AExecuteBatch(MAIN_KEYSPACE_ALIAS, batch)
}

func AExecuteBatch(alias string, batch *gocql.Batch) error {
	return Sessions[alias].ExecuteBatch(batch)
}
