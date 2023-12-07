package database

import (
	"Popkat/database/types"
	"Popkat/state"
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
)

func Cum() {
	Query := generateTableSQL(types.Service{}) + " " + generateTableSQL(types.User{})
	fmt.Println(Query)

	d, err := state.Pool.Exec(state.Context, Query)
	fmt.Println(d)

	if err != nil {
		panic(err)
	}
}

func generateTableSQL(dataModel interface{}) string {
	var columns []string

	t := reflect.TypeOf(dataModel)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}

		typeTag := field.Tag.Get("type")

		columnType := "TEXT" // default type if not specified in the tag

		if strings.Contains(typeTag, "varchar") {
			columnType = "VARCHAR(255)"
		} else if strings.Contains(typeTag, "int") {
			columnType = "INTEGER"
		} else if strings.Contains(typeTag, "boolean") {
			columnType = "BOOLEAN"
		} else if strings.Contains(typeTag, "timestamp") {
			columnType = "TIMESTAMP"
		} else if strings.Contains(typeTag, "jsonb") {
			columnType = "JSONB"
		} else if strings.Contains(typeTag, "float") {
			columnType = "REAL"
		} else if strings.Contains(typeTag, "double") {
			columnType = "DOUBLE PRECISION"
		} else if strings.Contains(typeTag, "uuid") {
			columnType = "UUID"
		}

		columns = append(columns, fmt.Sprintf("%s %s", tag, columnType))
	}

	tableName := pq.QuoteIdentifier(strings.ToLower(t.Name()) + "s")

	createTableSQL := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, strings.Join(columns, ", "))

	return createTableSQL
}
