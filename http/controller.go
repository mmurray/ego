package http


type Controller struct {
	currentAction *Action
}

type Action struct {
	Name string
	parameterNames []string
	resultNames []string
}