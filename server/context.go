package server

import "fmt"

type Context struct {
	params CommandParams
}

func (ctx *Context) Param(name string) (string, error) {
	value, exists := ctx.params[name]

	if !exists {
		return "", fmt.Errorf("no parameter: %s", name)
	}

	return value, nil
}
