package command

import (
	"fmt"
	m "simple-irc-server/message"
)

type Command[PT any] struct {
	Name         string
	Handler      Handler[PT]
	ParamsParser ParamsParser[PT]
}

type Handler[PT any] func(*Context[PT]) error

type Context[PT any] struct {
	Params *PT
}

type ParamsParser[PT any] = func(string) (*PT, error)

func (cmd *Command[P]) GetName() string {
	return cmd.Name
}

func (cmd *Command[PT]) Execute(msg *m.Message) error {
	params, err := cmd.parseParams(msg.Params())

	if err != nil {
		return err
	}

	context := &Context[PT]{params}

	return cmd.runHandler(context)
}

func (cmd *Command[PT]) parseParams(raw string) (*PT, error) {
	params, err := cmd.ParamsParser(raw)

	if err != nil {
		return nil, fmt.Errorf("error parsing params: %w", err)
	}

	return params, nil
}

func (cmd *Command[PT]) runHandler(ctx *Context[PT]) error {
	if err := cmd.Handler(ctx); err != nil {
		return fmt.Errorf("handler retuned error: %v", err)
	}

	return nil
}
