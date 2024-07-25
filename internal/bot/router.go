package bot

import (
	"github.com/yardnsm/leeches/internal/model"
	tele "gopkg.in/telebot.v3"
)

type commandHandler func(Context, tele.Context) error
type commandsMap map[interface{}]*Command

type createContextFunc func(tele.Context) Context

type Router struct {
	commands     commandsMap
	textCommands commandsMap

	createContext   createContextFunc
	defaultRestrict restrictFunc
}

func NewRouter() *Router {
	return &Router{
		commands:     make(commandsMap),
		textCommands: make(commandsMap),
	}
}

func (router *Router) CreateContext(cbc createContextFunc) *Router {
	router.createContext = cbc
	return router
}

func (router *Router) AddCommand(cmd *Command) *Router {
	router.commands[cmd.endpoint] = cmd
	return router
}

func (router *Router) AddCallback(cmd *Command) *Router {
	switch end := cmd.endpoint.(type) {
	case string:
		// Refer to telebot.v3@v3.1.2/callback.go, CallbackUnique for *InlineButton
		router.commands["\f"+end] = cmd
	case tele.CallbackEndpoint:
		router.commands[end] = cmd
	default:
		router.commands[end] = cmd
	}
	return router
}

func (router *Router) AddTextCommand(cmd *Command) *Router {
	router.textCommands[cmd.endpoint] = cmd
	return router
}

func (router *Router) DefaultRestrictUser(defaultRestrict restrictFunc) *Router {
	router.defaultRestrict = defaultRestrict
	return router
}

func (router *Router) Attach(b *tele.Bot) *Router {

	// Attach handlers
	for endpoint, cmd := range router.commands {
		b.Handle(
			endpoint,
			router.createCommandHandler(cmd),
			cmd.middlewares...,
		)
	}

	// Attach onText handler
	b.Handle(
		tele.OnText,
		router.handleTextCommand,
	)

	return router
}

func (router *Router) GetUserCommands(user model.User) (commands []tele.Command) {
	for endpoint, cmd := range router.commands {
		switch v := endpoint.(type) {
		case string:
			// Ignore if the command is telebot related
			if v[0] == '\a' ||  v[0] == '\f' {
				break
			}

			err := router.handleRestriction(cmd, &user)
			if err == nil {
				commands = append(commands, tele.Command{
					Text:        v,
					Description: cmd.description,
				})
			}
		default:
			break
		}
	}

	return commands
}

func (router *Router) handleRestriction(cmd *Command, user *model.User) error {
	restrict := cmd.restrict
	if restrict == nil {
		restrict = router.defaultRestrict
	}

	if restrict == nil {
		return nil
	}

	return restrict(user)
}

func (router *Router) createCommandHandler(cmd *Command) tele.HandlerFunc {
	return func(t tele.Context) error {
		context := router.createContext(t)

		// Restirctions
		if err := router.handleRestriction(cmd, context.CurrentUser); err != nil {
			return nil
		}

		// When a command will run, we'll clear the text command state for the chat. We assume that
		// text command state is set by commands, and therefore before running each command we want
		// a fresh text command state.
		context.ClearTextCommand()

		return cmd.handler(context, t)
	}
}

func (router *Router) handleTextCommand(t tele.Context) error {
	context := router.createContext(t)
	textCommand := context.GetTextCommand()

	cmd, exists := router.textCommands[textCommand.Command]

	if !exists {

		// Fallback to default OnText handler, if set
		cmd, exists = router.commands[tele.OnText]
		if !exists {
			return nil
		}
	}

	if err := router.handleRestriction(cmd, context.CurrentUser); err != nil {
		return nil
	}

	return cmd.handler(context, t)
}
