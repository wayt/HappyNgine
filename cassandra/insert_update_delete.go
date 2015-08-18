package cassandra

import (
	"fmt"
	"reflect"
	"strings"
)

func parseTag(tag string) (string, string) {

	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tag[idx+1:]
	}
	return tag, ""
}

func createInsertQuery(table string, data interface{}) (string, []interface{}, error) {

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", nil, errNilPointer
		}
		v = v.Elem()
	}

	columns := make([]string, 0)
	marks := make([]string, 0)
	values := make([]interface{}, 0)

	for i := 0; i < v.NumField(); i++ {

		if v.Field(i).CanInterface() {
			value := v.Field(i).Interface()
			name, _ := parseTag(v.Type().Field(i).Tag.Get(TAG_NAME))

			if name != "" && name != "-" {
				columns = append(columns, name)
				marks = append(marks, "?")
				values = append(values, value)
			}
		}
	}
	stmt := fmt.Sprintf(QUERY_INSERT, table, strings.Join(columns, QUERY_SEPARATOR), strings.Join(marks, QUERY_SEPARATOR))

	return stmt, values, nil
}

func Insert(table string, data interface{}, stmtExtras ...string) error {

	stmt, values, err := createInsertQuery(table, data)
	if err != nil {
		return err
	}

	// Useful to specify `USING TTL 42`
	stmt += " " + strings.Join(stmtExtras, "")

	if err := Query(stmt, values...).Exec(); err != nil {
		return err
	}

	return nil
}

func createUpdateQuery(table string, data interface{}) (string, []interface{}, error) {

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", nil, errNilPointer
		}
		v = v.Elem()
	}

	columns := make([]string, 0)
	where := make([]string, 0)
	values := make([]interface{}, 0)

	for i := 0; i < v.NumField(); i++ {

		if v.Field(i).CanInterface() {
			value := v.Field(i).Interface()
			name, option := parseTag(v.Type().Field(i).Tag.Get(TAG_NAME))

			if name != "" && name != "-" {
				if option == "key" {
					where = append(where, fmt.Sprintf(QUERY_ITEM, name))
					values = append(values, value)
				} else {
					columns = append([]string{fmt.Sprintf(QUERY_ITEM, name)}, columns...)
					values = append([]interface{}{value}, values...)
				}
			}
		}
	}
	stmt := fmt.Sprintf(QUERY_UPDATE, table, strings.Join(columns, QUERY_SEPARATOR), strings.Join(where, QUERY_AND))
	return stmt, values, nil
}

func Update(table string, data interface{}) error {

	stmt, values, err := createUpdateQuery(table, data)
	if err != nil {
		return err
	}
	if err := Query(stmt, values...).Exec(); err != nil {
		return err
	}

	return nil
}

func createDeleteQuery(table string, data interface{}) (string, []interface{}, error) {

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", nil, errNilPointer
		}
		v = v.Elem()
	}

	where := make([]string, 0)
	values := make([]interface{}, 0)

	for i := 0; i < v.NumField(); i++ {

		if v.Field(i).CanInterface() {
			value := v.Field(i).Interface()
			name, option := parseTag(v.Type().Field(i).Tag.Get(TAG_NAME))

			if name != "" && name != "-" && option == "key" {
				where = append(where, fmt.Sprintf(QUERY_ITEM, name))
				values = append(values, value)
			}
		}
	}
	stmt := fmt.Sprintf(QUERY_DELETE, table, strings.Join(where, QUERY_AND))
	return stmt, values, nil
}

func Delete(table string, data interface{}) error {

	stmt, values, err := createDeleteQuery(table, data)
	if err != nil {
		return err
	}
	if err := Query(stmt, values...).Exec(); err != nil {
		return err
	}

	return nil
}
