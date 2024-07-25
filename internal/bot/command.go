package bot

import (
	tele "gopkg.in/telebot.v3"
)

type Command struct {
	endpoint    interface{}
	description string
	restrict    restrictFunc
	handler     commandHandler
	middlewares []tele.MiddlewareFunc
}

func NewCommand(endpoint interface{}) *Command {
	return &Command{
		endpoint:    endpoint,
		middlewares: make([]tele.MiddlewareFunc, 0),
	}
}

func (cmd *Command) Description(description string) *Command {
	cmd.description = description
	return cmd
}

func (cmd *Command) RestrictUser(restrict restrictFunc) *Command {
	cmd.restrict = restrict
	return cmd
}

func (cmd *Command) Middleware(middleware tele.MiddlewareFunc) *Command {
	cmd.middlewares = append(cmd.middlewares, middleware)
	return cmd
}

func (cmd *Command) Handle(handler commandHandler) *Command {
	cmd.handler = handler
	return cmd
}
