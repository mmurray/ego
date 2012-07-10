package http

import (
	nhttp "net/http"
	"log"
	"strings"
)

type Action struct {
	Path string
	Method string
	ResponseTypes string
	Execute func(*Request) Result
	Package string
	View string
	request *Request
}

type Result interface{}

func (a *Action) Dispatch(w nhttp.ResponseWriter, httpReq *nhttp.Request, urlparams map[string]interface{}, rtype string) {
	w.Header().Set("Server", "ego")
	ctx := &RequestContext{
		Writer: &w,
	}
	req := NewRequest()
	req.Context = ctx
	req.Params = urlparams
	req.Parse(httpReq)
	result := a.Execute(req)
	resp, ok := result.(*Response)
	if !ok {
		resp = &Response{
			Context: map[string]interface{}{
				"result": result,
				"foo": "bar",
			},
		}
	}
	if (resp.StatusCode != 0) {
		w.WriteHeader(resp.StatusCode)
	}
	switch rtype {
	case "json":
		resp.WriteJSON(w)
	case "xml":
		resp.WriteXML(w)
	default:
		if _, isStr := result.(string); isStr {
			resp.WriteText(w)
		} else {
			if resp.View == "" {
				if a.View != "" {
					resp.View = a.View
				} else {
					log.Panic("Cannot respond with HTML, you must specify either Response.View or Action.View")
				}
			}
			resp.WriteHTML(w)
		}
	}
}

func (a *Action) HandlesMethod(method string) bool {
	m := a.Method
	if m == "" || m == "*" {
		return true; // by default actions handle all methods
	}
	return strings.ToUpper(method) == strings.ToUpper(m);
}

var NotFoundAction = &Action {
	Execute: func(*Request) Result {
		return NotFound
	},
}