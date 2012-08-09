package actions

import (
	"github.com/murz/ego/http"
	"github.com/murz/ego/ws"
	"reflect"
	"log"
)

type ActionManager struct {
	HTTPActions []*http.Action
	WSActions []*ws.Action
}

var am = &ActionManager{
	HTTPActions: make([]*http.Action, 0),
}

func Register(a interface{}) Action {
	switch act := reflect.ValueOf(a).Interface().(type) {
	case *http.Action:
		log.Printf("%v", act.Path)
		am.HTTPActions = append(am.HTTPActions, act)
		return act
	case *ws.Action:
		am.WSActions = append(am.WSActions, act)
		return act
	}
	return nil
}

func Count() int {
	log.Printf("%v", len(am.HTTPActions) + len(am.WSActions))
	return len(am.HTTPActions) + len(am.WSActions)
}

func HTTPActions() []*http.Action {
	return am.HTTPActions
}

func WSActions() []*ws.Action {
	return am.WSActions
}