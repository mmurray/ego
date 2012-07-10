package http

import (
	"regexp"
	"log"
	"strings"
	"fmt"
	nhttp "net/http"
)

type RouterNode struct {
	Value string
	ParamKeys []string
	Regexp *regexp.Regexp
	ChildNodes map[string]*RouterNode
	Action *Action
}

func (n *RouterNode) HasValue() bool {
	return n.Value != ""
}

func NewRouterNode() *RouterNode {
	return &RouterNode {
		ParamKeys: make([]string, 0),
		ChildNodes: make(map[string]*RouterNode),
		Value: "/",
	}
}

type Router struct {
	KeyRegexp *regexp.Regexp
	ParamRegexp *regexp.Regexp
	RootNode *RouterNode
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

func (r *Router) Keys(a *Action) []string {
	matches := strings.Split(a.Path, "/")
	if (len(matches) < 2) {
		log.Panicf("ego: \"%v\" is not a valid action path.", a.Path)
	}
	return matches
}

func (r *Router) Register(a *Action) {
	keys := r.Keys(a)
	curNode := r.RootNode
	if a.Method == "" {
		a.Method = "*"
	}

	if curNode.ChildNodes[a.Method] != nil {
		curNode = curNode.ChildNodes[a.Method]
	} else {
		node := NewRouterNode()
		node.Value = a.Method
		curNode.ChildNodes[a.Method] = node
		curNode = node
	}

	keys = keys[1:len(keys)]
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
	curNode.Action = a
}

func (r *Router) Lookup(path string, method string) (*Action, map[string]interface{}, bool) {
	tokens := strings.Split(path, "/")
	tokens = tokens[1:len(tokens)] // paths look like "/foo" so the first token will always be ""
	curNode := r.RootNode.ChildNodes[method]
	if curNode == nil {
		return nil, nil, false
	}
	match := false
	globparams := make(map[string]interface{})
	for i, token := range tokens {
		found := false
		log.Printf("-- level%v: %v",i, token)
		if token == "" && i == len(tokens) - 1 {
			match = true
			break;
		}
		log.Printf("%v", curNode.ChildNodes)
		if (token == curNode.Value) {
			if  i == len(tokens) - 1 {
				match = true
				break;
			}
		}
		for _, node := range curNode.ChildNodes {
			log.Printf("nv: %v", node.Value)
			log.Printf("t: %v", token)
			if node.Value == token {
				curNode = node
				found = true
				if  i == len(tokens) - 1 {
					match = true
					break;
				}
			} else if node.Regexp != nil {
				log.Printf("rxp: %v", node.Regexp)
				matches := node.Regexp.FindAllStringSubmatch(token, -1)
				log.Printf("matches; %v", matches)
				if len(matches) >= 1 {
					curNode = node
					found = true
					var params []string
					if len(matches[0]) >= 2 {
						params = matches[0][1:len(matches[0])]
					} else {
						params = matches[0]
					}
					log.Printf("%v", curNode.ParamKeys)
					log.Printf("%v vs %v", len(params), len(curNode.ParamKeys))
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
	if (match && curNode.Action != nil) {
		return curNode.Action, globparams, true
	}
	return nil, nil, false
	// for len(tokens) != 0 {

	// }
}

func (r *Router) ActionDispatchHandler() nhttp.HandlerFunc {
	return func(w nhttp.ResponseWriter, httpReq *nhttp.Request) {
		reqType := ""
		a, p, found := r.Lookup(httpReq.URL.Path, httpReq.Method)
		if !found {
			// try the wildcard tree
			a, p, found = r.Lookup(httpReq.URL.Path, "*")
			if !found {
				NotFoundAction.Dispatch(w, httpReq, nil, reqType)
				return;
			}
		}
		a.Dispatch(w, httpReq, p, reqType)
	}
}