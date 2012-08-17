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
	Layout string
	Text string
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
				"Result": result,
			},
		}
	} else {
		ctx, ok := resp.Context.(map[string]interface{})
		if !ok {
			ctx = make(map[string]interface{})
			ctx["Result"] = resp.Context
			resp.Context = ctx
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
		if txt, isStr := result.(string); isStr {
			resp = &Response{
				Context: txt,
			}
			resp.WriteText(w)
		} else {
			log.Printf("resp view: %v", resp.View)
			if resp.View == "" {
				if a.View != "" {
					resp.View = a.View
				} else if resp.Text != "" {
					resp.Context = resp.Text
					resp.WriteText(w)
					return
				} else {
					log.Panic("Cannot respond with HTML, you must specify either Response.View or Action.View")
				}
			}
			if resp.Layout == "" {
				if a.Layout != "" {
					resp.Layout = a.Layout
				} else {
					resp.Layout = "application.html"
				}
			}
			if mp, isMap := resp.Context.(map[string]interface{}); isMap {
				log.Printf("ctx: %v", mp["Result"])
			} else {
				log.Printf("not a map")
			}
			log.Printf("ctx: %v", resp.Context)
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