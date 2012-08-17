package db

import (
	"fmt"
	"reflect"
	"log"
	"strings"
)

type Query struct {
	model interface{}
	val reflect.Value
	keyField string
	table string
	where string
	set string
	params []interface{}
	paramWrapper string
	nextParamIndex int
}

func NewQuery(model interface{}) *Query {
	val := reflect.Indirect(reflect.ValueOf(model))
	keyField := "Id"
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		if field.Tag == "key" {
			keyField = field.Name
		}
	}
	query := &Query{
		model: model,
		keyField: keyField,
		val: val,
		table: val.Type().Name(),
		params: make([]interface{}, 0),
		nextParamIndex: 0,
	}
	switch con.driver {
	case "postgres":
		query.paramWrapper = "\""
	default:
		query.paramWrapper = "`"
	}
	return query
}

func (q *Query) placeholder() string {
	if con.driver == "postgres" {
		q.nextParamIndex++
		return fmt.Sprintf("$%d", q.nextParamIndex)
	}
	return "?"
}

func (q *Query) WhereKeyEquals(key interface{}) *Query {
	return q.Where(fmt.Sprintf("where %s%s%s=%v", q.paramWrapper, q.keyField, q.paramWrapper, q.placeholder()), key)
}

func (q *Query) Where(where string, args ...interface{}) *Query {
	q.where = where
	q.params = append(q.params, args...)
	return q
}

func (q *Query) ToSelect() string {
	return fmt.Sprintf("SELECT * from %s %s", q.table, q.where)
}

func (q *Query) ToUpdate() string {
	return fmt.Sprintf("UPDATE %s set %s %s", q.table, q.set, q.where)
}

func (q *Query) Fetch() interface{} {
	sql := q.ToSelect()
	log.Print(sql)
	ps, err := con.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer ps.Close() 
	rows, err := ps.Query(q.params...)
	if err != nil {
		panic(err)
	}
	gotResults := rows.Next()
	if !gotResults {
		return nil
	}
	cols, _ := rows.Columns()
	pointers := make([]interface{}, len(cols)) 
	for key, name := range cols {
		field := q.val.FieldByName(strings.Title(name))
		pointers[key] = field.Addr().Interface()
	}
	rows.Scan(pointers...)
	return q.model
}

func (q *Query) FetchAll() []interface{} {
	sql := q.ToSelect()
	log.Print(sql)
	ps, err := con.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer ps.Close()
	rows, err := ps.Query(q.params...)
	if err != nil {
		panic(err)
	}
	result := make([]interface{}, 0)
	for rows.Next() {
		val := reflect.New(q.val.Type()).Elem()
		cols, _ := rows.Columns()
		pointers := make([]interface{}, len(cols)) 
		for key, name := range cols {
			field := val.FieldByName(strings.Title(name))
			pointers[key] = field.Addr().Interface()
		}
		rows.Scan(pointers...)
		result = append(result, val.Interface())
	}
	return result
}

func (q *Query) Save() *DBResult {
	val := reflect.Indirect(reflect.ValueOf(q.model))
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		fieldval := val.Field(i)
		if field.Name == q.keyField {
			q.WhereKeyEquals(fieldval.Interface())
		} else {
			if q.set != "" {
				q.set += ", "
			}
			valstring := ""
			switch f := fieldval.Interface().(type) {
			case string:
				valstring = fmt.Sprintf("'%v'", f)
			case bool:
				if f {
					valstring = "true"
				} else {
					valstring = "false"
				}
			}
			q.set += fmt.Sprintf("%v%v%v=%v", q.paramWrapper, field.Name, q.paramWrapper, valstring)
		}
		
	}
	log.Printf("sql: %v", q.ToUpdate())
	ps, err := con.db.Prepare(q.ToUpdate())
	if err != nil {
		panic(err)
	}
	defer ps.Close()
	ps.Exec(q.params...)
	return &DBResult {
		Success: true,
		Message: "Your record was successfully saved",
	}
} 