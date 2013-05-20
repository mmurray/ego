package http

import (
	"github.com/murz/ego/tmpl"
	nhttp "net/http"
	"io"
	"encoding/json"
	"encoding/xml"
	"log"
)

type Response struct {
	View string
	Layout string
	Text string
	StatusCode int
	Context map[string]interface{}
	Type string
}

func (r *Response) WriteHTML(w nhttp.ResponseWriter) {
	if r.View == "" {
		log.Panic("Attempted to call Response.WriteHTML when Response.View was not set")
	}
	log.Printf("ctx: %v", r.Context)
	tmpl.Render(w, r.View, r.Context)
	
	// tmpl.Render(w, r.View, r.Layout, r.Context)
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
	io.WriteString(w, r.Text)
}

var NotFound = &Response{
	StatusCode: 404,
	View: "/errors/404.html",
	Layout: "none",
}

var NotImplemented = &Response{
	StatusCode: 501,
	View: "/errors/501.html",
	Layout: "none",
}