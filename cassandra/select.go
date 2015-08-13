package cassandra

import (
	"fmt"
	"github.com/gocql/gocql"
	"reflect"
	"strings"
)

const (
	ASC = iota
	DESC
)

type selectQuery struct {
	selec      string
	where      string
	orderBy    string
	limit      string
	parameters []interface{}
}

func Select(table string) *selectQuery {

	query := &selectQuery{
		QUERY_SELECT + table,
		"",
		"",
		"",
		make([]interface{}, 0),
	}

	return query
}

func (q *selectQuery) Where(condition string, values ...interface{}) *selectQuery {

	if len(q.where) == 0 {
		q.where = " WHERE " + condition
	} else {
		q.where = q.where + QUERY_AND + condition
	}

	q.parameters = append(q.parameters, values...)

	return q
}

func (q *selectQuery) OrderBy(param string, order int) *selectQuery {

	q.orderBy = " ORDER BY " + param

	if order == ASC {
		q.orderBy = q.orderBy + " ASC"
	} else {
		q.orderBy = q.orderBy + " DESC"
	}

	return q
}

func (q *selectQuery) Limit(param int) *selectQuery {

	q.limit = fmt.Sprintf(" LIMIT %d", param)

	return q
}

func (q *selectQuery) GetQuery(columns []string) string {

	str := fmt.Sprintf(q.selec+q.where+q.orderBy+q.limit, strings.Join(columns, QUERY_SEPARATOR))
	fmt.Println(str)

	return str
}

func (q *selectQuery) Scan(object interface{}) (bool, error) {

	columns := make([]string, 0)
	values := make([]interface{}, 0)

	v := reflect.ValueOf(object)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false, errNotAPointer
		}
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i).Addr().Interface()
		tag, _ := parseTag(v.Type().Field(i).Tag.Get(TAG_NAME))

		if tag != "" && tag != "-" {
			columns = append(columns, tag)
			values = append(values, f)
		}
	}

	if err := Query(q.GetQuery(columns), q.parameters...).Scan(values...); err != nil {
		if err == gocql.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (q *selectQuery) ScanAll(objects interface{}) (bool, error) {

	v := reflect.ValueOf(objects)
	if v.Kind() != reflect.Ptr {
		return false, errNotAPointer
	}
	if v.IsNil() {
		return false, errNilPointer
	}
	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return false, errNotASlice
	}

	elmType := v.Type().Elem()
	if elmType.Kind() != reflect.Ptr {
		return false, errNotAPointer
	}
	elmType = elmType.Elem()
	columns := make([]string, 0)
	for i := 0; i < elmType.NumField(); i++ {
		tag, _ := parseTag(elmType.Field(i).Tag.Get(TAG_NAME))
		if tag != "" && tag != "-" {
			columns = append(columns, tag)
		}
	}
	iter := Query(q.GetQuery(columns), q.parameters...).Iter()

	for {
		object := reflect.New(elmType).Elem()
		values := make([]interface{}, 0)

		for i := 0; i < object.NumField(); i++ {
			f := object.Field(i).Addr().Interface()
			tag, _ := parseTag(object.Type().Field(i).Tag.Get(TAG_NAME))
			if tag != "" && tag != "-" {
				values = append(values, f)
			}
		}
		if !iter.Scan(values...) {
			break
		}
		v.Set(reflect.Append(v, object.Addr()))
	}

	if err := iter.Close(); err != nil {
		return false, err
	}

	return true, nil
}
