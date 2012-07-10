package ws

import (
	"fmt"
	"socketio"
)

type Action struct {
	Event string
	Execute func(args ...interface{}) *Result
}

type Result interface{}

func (a *Action) Dispatch(msg *socketio.Message) {
	args := [...]interface{}{

	}
    name, _ := msg.Event()
    msg.ReadArguments(args)
    fmt.Printf("this is an event with name %s.", name)
    a.Execute(args)
}