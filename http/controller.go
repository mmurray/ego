package http


type Controller struct {
	currentAction *Action
}

type Action struct {
	Name string
	parameterNames []string
	resultNames []string
}

// A helper function to quickly create responses from a set of values.
func (c *Controller) Result(values ...interface{}) *Response {
	return &Response{
		Context: values[0],
	}
}