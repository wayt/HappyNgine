package cassandra

import (
	"errors"
)

var TAG_NAME = "cassandra"

var QUERY_INSERT = "INSERT INTO %s (%s) VALUES(%s)"
var QUERY_UPDATE = "UPDATE %s SET %s WHERE %s"
var QUERY_DELETE = "DELETE FROM %s WHERE %s"
var QUERY_SELECT = "SELECT %s FROM "
var QUERY_ITEM = "%s = ?"
var QUERY_SEPARATOR = ", "
var QUERY_AND = " AND "

var errNilPointer = errors.New("happyngine/cassandra: The pointer is nil")
var errNotAPointer = errors.New("happyngine/cassandra: It's not a pointer")
var errNotASlice = errors.New("happyngine/cassandra: It's not a slice")
