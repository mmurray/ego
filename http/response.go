package http

import (
	"ego/tmpl"
	nhttp "net/http"
	"io"
	"encoding/json"
	"encoding/xml"
	"log"
)

type Response struct {
	View string
	StatusCode int
	Context interface{}
	Type string
}

func (r *Response) WriteHTML(w nhttp.ResponseWriter) {
	if r.View == "" {
		log.Panic("Attempted to call Response.WriteHTML when Response.View was not set")
	}
	io.WriteString(w, tmpl.Render(r.View, r.Context))
}

func (r *Response) WriteJSON(w nhttp.ResponseWriter) {
	b, err := json.Marshal(r.Context)
	if err != nil {
		log.Panic("Error marshalling JSON response")
	}
	w.Write(b)
}

func (r *Response) WriteXML(w nhttp.ResponseWriter) {
	b, err := xml.Marshal(r.Context)
	if err != nil {
		log.Print("Error marshalling XML response")
		panic(err)
	}
	w.Write(b)
}

func (r *Response) WriteText(w nhttp.ResponseWriter) {
	str, ok := r.Context.(string)
	if ok {
		io.WriteString(w, str)
	} else {
		log.Panic("Attempted to call Response.WriteText when Response.Context was not a string")
	}
}

var NotFound = &Response{
	StatusCode: 404,
	View: "/errors/404.html",
}