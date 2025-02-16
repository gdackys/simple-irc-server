package server

import (
	m "simple-irc-server/message"
)

type Command struct {
	Name         CommandName
	ParamsParser ParamsParser
	Handler      CommandHandler
}

type CommandName = string
type ParamsParser = func(string) (CommandParams, error)
type CommandHandler func(*Context) error

type CommandParams = map[string]string

func (cmd *Command) execute(message *m.Message) error {
	params, err := cmd.parseParams(message.Params())

	if err != nil {
		return err
	}

	context := &Context{params}

	return cmd.runHandler(context)
}

func (cmd *Command) parseParams(rawParams string) (CommandParams, error) {
	return cmd.ParamsParser(rawParams)
}

func (cmd *Command) runHandler(ctx *Context) error {
	return cmd.Handler(ctx)
}
