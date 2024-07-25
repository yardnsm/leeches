package bot

import (
	"fmt"

	"github.com/yardnsm/gohever"
	"github.com/yardnsm/leeches/internal/model"

	tele "gopkg.in/telebot.v3"
)

type TextCommandForChat struct {
	Command interface{}
	State   interface{}
}

// This will store a global map in memory for text command for each chat
var (
	textCommandForChats map[int64]TextCommandForChat
)

func init() {
	textCommandForChats = make(map[int64]TextCommandForChat)
}

// bot.Context is a wrapper around tele.Context, which have access to the hever API and the DBs
// repository, plus more neat stuff
type Context struct {
	tc tele.Context

	CurrentUser    *model.User
	Users          *model.UsersRepository
	ChargeRequests *model.ChargeRequestsRepository
	Hever          *gohever.Client
}

func NewContext(tc tele.Context) (context Context) {
	context.tc = tc
	return context
}

// Returns the underlying bot instance
func (c *Context) Bot() *tele.Bot {
	return c.tc.Bot()
}

// Get the state of the current text command set for the current chat. The text command mechanism
// is bound to a specific chat, by it's ID
func (c *Context) GetTextCommand() TextCommandForChat {
	chatID := c.tc.Chat().ID
	return textCommandForChats[chatID]
}

// Set the state for the next text command
func (c *Context) SetTextCommand(command interface{}, state interface{}) {
	chatID := c.tc.Chat().ID
	textCommandForChats[chatID] = TextCommandForChat{
		Command: command,
		State:   state,
	}
}

// Clear the current text command.
// This is fired automatically when a "regular" (/) command is fired, in order to keep the text
// command state clean after each command.
func (c *Context) ClearTextCommand() {
	chatID := c.tc.Chat().ID
	delete(textCommandForChats, chatID)
}

func (c *Context) SendAndCloseKeyboard(what interface{}) error {
	return c.tc.Send(what, &tele.ReplyMarkup{RemoveKeyboard: true})
}

func (c *Context) SendError(err error) error {
	return c.tc.Send(fmt.Sprintf("*Error:*\n`%s`", err), tele.ModeMarkdownV2)
}

// Send a message to the current recipient, then return the new message
func (c *Context) SendEditable(what interface{}, opts ...interface{}) (*tele.Message, error) {
	return c.tc.Bot().Send(c.tc.Recipient(), what, opts...)
}

// A shorthand to tele.Bot.Edit that returns only an error
func (c *Context) Edit(msg tele.Editable, what interface{}, opts ...interface{}) error {
	_, err := c.tc.Bot().Edit(msg, what, opts...)
	return err
}
