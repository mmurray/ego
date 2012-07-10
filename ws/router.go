package ws

import (
	"log"
	"socketio"
	"fmt"
)

type Router struct {
	Actions map[string]*Action
}

func NewRouter() *Router {
	return &Router{
		Actions: make(map[string]*Action, 0),
	}
}

func (r *Router) Register(a *Action) {
	r.Actions[a.Event] = a
}

func (r *Router) Lookup(evt string) (*Action, bool) {
	action := r.Actions[evt]
	if action == nil {
		return nil, false
	}
	return action, true
}

func (r *Router) ActionDispatchHandler() func(c *socketio.Conn) {
	return func(c *socketio.Conn) {
		var msg socketio.Message
	    for {
	        if err := c.Receive(&msg); err != nil {
	            return
	        }
	        switch msg.Type() {
	        case socketio.MessageJSON, socketio.MessageText:
	            fmt.Println("this is an ordinary message with payload: ", msg.String())
	        case socketio.MessageEvent:
	            name, _ := msg.Event()
	            if act, ok := r.Lookup(name); ok {
	            	act.Dispatch(&msg)
	            }
	        }
	    }
		log.Printf("yay! %v", c)
	}
}