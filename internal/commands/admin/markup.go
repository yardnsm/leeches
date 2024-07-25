package admin

import (
	"fmt"

	"github.com/yardnsm/leeches/internal/bot"

	tele "gopkg.in/telebot.v3"
)

func createMarkup(c bot.Context, isAdmin bool) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	allUsers, _ := c.Users.GetAll()
	var menuRows []tele.Row

	for _, user := range allUsers {

		// Skip current user
		if user.ID == c.CurrentUser.ID {
			continue
		}

		if user.IsAdmin == isAdmin {
			menuRows = append(menuRows, menu.Row(menu.Text(
				fmt.Sprintf("%s (%d)", user.DisplayName, user.ID),
			)))
		}
	}

	menu.Reply(menuRows...)
	return menu
}
