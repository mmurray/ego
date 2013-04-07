package http

import (
	"reflect"
	"fmt"
	"log"
	"time"
	nhttp "net/http"
)

// An http result can be anything, the dispatcher will try to figure out
// what to do with the result based on it's type.
type Result interface {}

// Metadata stored about the controllers at startup time.
type ControllerMetadata struct {
	Type reflect.Type
}

var controllerRegistry = make(map[string]*ControllerMetadata)

// Bind action key in the form of "Controller.Action" to a controller type.
func RegisterAction(key string, t reflect.Type) {
	// TODO: handle error case
	controllerRegistry[key] = &ControllerMetadata{
		t,
	}
}

// Returns a net/http handler function for dispatching ego actions using the given router.
func ActionDispatchHandler(r *Router) nhttp.HandlerFunc {
	return func(w nhttp.ResponseWriter, httpReq *nhttp.Request) {
		startTime := time.Now().UnixNano() / 1000000
		// reqType := ""
		route, _, found := r.Lookup(httpReq.URL.Path, httpReq.Method)
		if !found {
			// try the wildcard tree
			route, _, found = r.Lookup(httpReq.URL.Path, "*")
			if !found {
				log.Printf("not found");
				// NotFoundAction.Dispatch(w, httpReq, nil, reqType)
				return;
			}
		}
		ctrlType := controllerRegistry[fmt.Sprintf("%s.%s", route.ControllerName, route.ActionName)].Type
		ctrlVal := reflect.New(ctrlType)
		cfgMethod := ctrlVal.MethodByName("Configure")
		if cfgMethod.IsValid() {
			cfgMethod.Call([]reflect.Value{})
		}
		method := ctrlVal.MethodByName(route.ActionName)
		resultVal := method.Call([]reflect.Value{})[0]
		result := resultVal.Interface()
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
		log.Printf("%s", resp)
		resp.WriteJSON(w)
		log.Printf("%s %s -> %s.%s (%dms)",
			route.Path.Method,
			route.Path.Value,
			route.ControllerName,
			route.ActionName,
			(time.Now().UnixNano() / 1000000) - startTime);
	}
}