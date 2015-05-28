package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/wayt/happyngine/env"
	"github.com/wayt/happyngine/log"
	"strings"
)

var ClusterCfg *gocql.ClusterConfig
var Session *gocql.Session

func init() {

	hostsString := env.Get("HAPPY_CASSANDRA_HOSTS")
	hosts := strings.Split(hostsString, ",")

	ClusterCfg = gocql.NewCluster(hosts...)
	ClusterCfg.Keyspace = env.Get("HAPPY_CASSANDRA_KEYSPACE")

	var err error
	Session, err = ClusterCfg.CreateSession()
	if err != nil {
		log.Criticalln(err)
	}
}

func TimeUUID() gocql.UUID {
	return gocql.TimeUUID()
}

func ParseUUID(input string) (gocql.UUID, error) {
	return gocql.ParseUUID(input)
}

func Query(stmt string, values ...interface{}) *gocql.Query {
	return Session.Query(stmt, values...)
}

func ExecuteBatch(batch *gocql.Batch) error {
	return Session.ExecuteBatch(batch)
}
