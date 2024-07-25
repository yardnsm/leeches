package bot

import (
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func AllowOnlyPrivateChatsMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(t tele.Context) error {
			if t.Chat().Type != tele.ChatPrivate {
				return nil
			}

			return next(t)
		}
	}
}

// TODO handle panics, I want ti catch those as well. There is a premade middleware called
// "middleware.Recover".
// TODO there are a kit of ways we can improve this, such as using the builtin
// tele.Settings.OnError.
// TODO send the error to all the admins for debugging purposes
func SendErrorsToUsersChatMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(t tele.Context) error {
			err := next(t)

			if err != nil && strings.Index(err.Error(), "telebot:") != 0 {
				err = t.Send(fmt.Sprintf("*Error:*\n`%s`", err), tele.ModeMarkdownV2)
			}

			return err
		}
	}
}
