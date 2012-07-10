package db

import (
	"ego/cfg"
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
)

type Connection struct {
	db *sql.DB
	driver string
}

var con = &Connection{}

func Connect(c *cfg.ConfigMap) {
	conf := (*c)
	driver := fmt.Sprintf("%s", conf["driver"])
	db, err := sql.Open(driver,
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", conf["user"], conf["password"], conf["name"]))
	if err != nil {
		panic(err)
	}
	con.db = db
	con.driver = driver
}

type DBResult struct {
	Success bool
	Message string
}

func GetById(model interface{}, id interface{}) interface{} {
	return NewQuery(model).WhereKeyEquals(id).Fetch()
}

// func GetBySQL(sql string, args ...interface{}) {
// 	return NewQuery(model).GetBySQL(sql, args)
// }

func GetAll(model interface{}) []interface{} {
	return NewQuery(model).FetchAll()
}

func Get(model interface{}) *Query {
	return NewQuery(model)
}

func Save(model interface{}) *DBResult {
	return NewQuery(model).Save()
}



// func (q *Query)  *Query