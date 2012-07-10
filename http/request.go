package http

import (
	nhttp "net/http"
	"reflect"
	"log"
)

type Request struct {
	Params map[string]interface{}
	Cookies map[string]string
	Context *RequestContext
}

type RequestContext struct {
	Writer *nhttp.ResponseWriter
}

func NewRequest() *Request {
	return &Request {
		Params: make(map[string]interface{}),
		Cookies: make(map[string]string),
	}
}

func (r *Request) Parse(httpreq *nhttp.Request) {
	httpreq.ParseForm()
	for key, val := range httpreq.Form {
		if len(val) == 1 {
			r.Params[key] = val[0]
		} else {
			r.Params[key] = val
		}
	}
}

func (r *Request) Populate(o interface{}) {
	val := reflect.Indirect(reflect.ValueOf(o))
	log.Printf("type: %v", val.Type().Name())
	for i := 0; i < val.Type().NumField(); i++ {
		structField := val.Type().Field(i)
		field := val.Field(i)
		log.Printf("field: %v", structField.Name)
		log.Printf("switching")
		param := r.Params[structField.Name]
		switch field.Kind() {
		case reflect.Bool:
			log.Printf("bool")
			if (param == nil) {
				field.Set(reflect.ValueOf(false))
			}
			if s, ok := param.(string); ok {
				if s == "on" {
					field.Set(reflect.ValueOf(true))
				}
			}
		case reflect.String:
			field.Set(reflect.ValueOf(param))
		}
	}
}