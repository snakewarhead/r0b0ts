package models

import (
	"database/sql"
	"reflect"

	// self contain
	_ "github.com/mattn/go-sqlite3"
)

const (
	// dbPath = "./resources/chain_data.db"
	dbPath = "./resources/chain_data_local.db"
)

var (
	db *sql.DB
)

func init() {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
}

// list the pointer of feilds of the struct of the model struct, usually use in the row.scan
// note that the model must be a pointer to a struct
func dbColumns(model interface{}) []interface{} {
	v := reflect.ValueOf(model).Elem()

	num := v.NumField()
	fields := make([]interface{}, num, num)
	for i := 0; i < num; i++ {
		fields[i] = v.Field(i).Addr().Interface()
	}

	return fields
}
