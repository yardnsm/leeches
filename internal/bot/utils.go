package bot

import (
	"strconv"

	tele "gopkg.in/telebot.v3"
)

// A command handler with state, useful for text commands.
// When an error is returned, the text command will be cleared.
// Returns a pointer to the new state since I want to accept nil.
type statefulCommandHandler[T any] func(Context, tele.Context, T) (interface{}, *T, error)

func SetCommandsForChat(b *tele.Bot, commands []tele.Command, chatID int64) {
	scope := tele.CommandScope{
		Type:   tele.CommandScopeChat,
		ChatID: chatID,
	}

	b.SetCommands(commands, scope)
}

func CreateStoredMessage(msg *tele.Message) tele.StoredMessage {
	return tele.StoredMessage{
		MessageID: strconv.Itoa(msg.ID),
		ChatID:    msg.Chat.ID,
	}
}

func EditableToStoredMessage(editable tele.Editable) tele.StoredMessage {
	messageId, chatId := editable.MessageSig()

	return tele.StoredMessage{
		MessageID: messageId,
		ChatID:    chatId,
	}
}

func CreateStatefulHandler[T any](handler statefulCommandHandler[T]) commandHandler {
	return func(c Context, t tele.Context) error {
		// Pls note that this will panic when type conversion fails.
		prevState := c.GetTextCommand().State.(T)

		nextCommand, nextState, err := handler(c, t, prevState)

		// Re-set the current command when the next is not defined, but there is a state
		// prevCommand := c.GetTextCommand().Command
		// if nextCommand == nil && nextState != nil {
		// 	c.SetTextCommand(prevCommand, *nextState)
		// 	return err
		// }

		if nextCommand != nil && err == nil {
			c.SetTextCommand(nextCommand, *nextState)
		} else {
			c.ClearTextCommand()
		}

		return err
	}
}
