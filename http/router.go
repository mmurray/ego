package http

import (
	"regexp"
	"log"
	"strings"
	"fmt"
)

// Routes in ego are stored in a tree data structure for quick lookup.
// RouterNode represents a single node in that tree.
type RouterNode struct {
	Value string
	ParamKeys []string
	Regexp *regexp.Regexp
	ChildNodes map[string]*RouterNode
	Route *Route
}

// Represents a router, which contains a tree of RouterNodes.
type Router struct {
	KeyRegexp *regexp.Regexp
	ParamRegexp *regexp.Regexp
	RootNode *RouterNode
}

// Represents a binding from a path to a controller action.
type Route struct {
	ControllerName string
	ActionName string
	Path Path
}

// Represents an http action path.
type Path struct {
	Value string
	Method string
}

// A helper for easily registering routes.
type RouteBuilder struct {
	path *Path
}

// Builds a wildcard path with the given value.
func Match(value string) *RouteBuilder {
	return &RouteBuilder{
		path: &Path{
			Value: value,
			Method: "*",
		},
	}
}

// Builds a GET path with the given value.
func Get(value string) *RouteBuilder {
	return &RouteBuilder{
		path: &Path{
			Value: value,
			Method: "GET",
		},
	}
}

func NewRouterNode() *RouterNode {
	return &RouterNode {
		ParamKeys: make([]string, 0),
		ChildNodes: make(map[string]*RouterNode),
		Value: "/",
	}
}

func NewRouter() *Router {
	r := &Router{
		RootNode: NewRouterNode(),
	}
	re, _ := regexp.Compile("([a-zA-Z?]{0,})")
	r.KeyRegexp = re
	pre, _ := regexp.Compile("{([a-zA-Z:]+)}")
	r.ParamRegexp = pre
	return r
}

var router = NewRouter()

func GetDefaultRouter() *Router {
	return router
}

// Binds the path to a route and registers it as a node in the router tree.
func (rb RouteBuilder) To(value string) {
	pieces := strings.Split(value, ".")
	if (len(pieces) != 2) {
		log.Panicf("ego: \"%v\" is not a valid controller action.", value)
	}
	router.Register(&Route{
		Path: *rb.path,
		ControllerName: pieces[0],
		ActionName: pieces[1],
	})
}

func (r *Router) Keys(p *Path) []string {
	path := p.Value
	if (path[0:1] == "/") {
		path = path[1:len(path)] // shave off the first "/" if there is one.
	}
	matches := strings.Split(path, "/")
	return matches
}

func (r *Router) Register(route *Route) {
	p := &route.Path
	keys := r.Keys(p)
	curNode := r.RootNode

	if curNode.ChildNodes[p.Method] != nil {
		curNode = curNode.ChildNodes[p.Method]
	} else {
		node := NewRouterNode()
		node.Value = p.Method
		curNode.ChildNodes[p.Method] = node
		curNode = node
	}

	for len(keys) != 0 {
		key := keys[0]
		keys = keys[1:len(keys)]
		if key == "" {
			break;
		}
		found := false
		for nodeKey, node := range curNode.ChildNodes {
			if key == nodeKey {
				curNode = node
				found = true
			}
		}
		if !found {
			log.Printf("creating node: %v", key)
			node := NewRouterNode()
			node.Value = key
			curNode.ChildNodes[key] = node
			curNode = node
			params := r.ParamRegexp.FindAllStringSubmatch(curNode.Value, -1)
			for _, param := range params {
				curNode.ParamKeys = append(curNode.ParamKeys, param[1])
			}
			log.Printf("lparams: %v", len(params))
			if len(params) > 0 {
				exps := r.ParamRegexp.ReplaceAllLiteralString(curNode.Value, "([a-zA-Z0-9]+)")
				log.Print(exps)
				log.Print(fmt.Sprintf("%v%v%v", "^", exps, "$"))
				exp, _ := regexp.Compile(fmt.Sprintf("%v%v%v", "^", exps, "$"))
				curNode.Regexp = exp
			}
		}
	}
	curNode.Route = route
}

func (r *Router) Lookup(path string, method string) (*Route, map[string]interface{}, bool) {
	curNode := r.RootNode.ChildNodes[method]
	if curNode == nil {
		return nil, nil, false
	}
	if (path[0:1] == "/") {
		path = path[1:len(path)] // shave off the first "/" if there is one.
	}
	tokens := strings.Split(path, "/")
	match := false
	globparams := make(map[string]interface{})
	for i, token := range tokens {
		found := false
		if token == "" && i == len(tokens) - 1 {
			match = true
			break;
		}
		if (token == curNode.Value) {
			if  i == len(tokens) - 1 {
				match = true
				break;
			}
		}
		for _, node := range curNode.ChildNodes {
			if node.Value == token {
				curNode = node
				found = true
				if  i == len(tokens) - 1 {
					match = true
					break;
				}
			} else if node.Regexp != nil {
				matches := node.Regexp.FindAllStringSubmatch(token, -1)
				if len(matches) >= 1 {
					curNode = node
					found = true
					var params []string
					if len(matches[0]) >= 2 {
						params = matches[0][1:len(matches[0])]
					} else {
						params = matches[0]
					}
					for i, key := range curNode.ParamKeys {
						globparams[key] = params[i]
					}
					if  i == len(tokens) - 1 {
						match = true
						break;
					}
				}
			}
		}
		if !found {
			break;
		}
	}
	if (match && curNode.Route != nil) {
		return curNode.Route, globparams, true
	}
	return nil, nil, false
	// for len(tokens) != 0 {

	// }
}