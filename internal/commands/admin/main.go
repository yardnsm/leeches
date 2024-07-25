package admin

import (
	"fmt"
	"strconv"

	"github.com/yardnsm/leeches/internal/bot"

	tele "gopkg.in/telebot.v3"
)

func setAdminOnSelection(c bot.Context, selection string, isAdmin bool) error {
	defer c.ClearTextCommand()
	extracted := selectionRegexp.FindStringSubmatch(selection)

	if len(extracted) == 0 {
		return c.SendAndCloseKeyboard("An invalid selection was made.")
	}

	userID, _ := strconv.ParseUint(extracted[1], 10, 0)

	// Do not allow the current user to set / unset himself
	if uint(userID) == c.CurrentUser.ID {
		return c.SendAndCloseKeyboard("You already an admin, buddy.")
	}

	user, err := c.Users.SetUserAdmin(uint(userID), isAdmin)
	if err != nil {
		c.SendAndCloseKeyboard("An error occured while updating the user")
		return err
	}

	fmtString := "%s is now an admin."
	if !isAdmin {
		fmtString = "%s is not an admin."
	}

	return c.SendAndCloseKeyboard(fmt.Sprintf(fmtString, user.DisplayName))
}

var setAdminCommand = bot.NewCommand(setAdminEndpoint).
	Description("Set an admin").
	RestrictUser(bot.RestrictAdmin).
	Handle(func(c bot.Context, t tele.Context) error {
		menu := createMarkup(c, false)

		t.Send("Please select a user to be set an admin.", menu)
		c.SetTextCommand(selectAdminToSetEndpoint, nil)

		return nil
	})

var removeAdminCommand = bot.NewCommand(removeAdminEndpoint).
	Description("Remove an admin").
	RestrictUser(bot.RestrictAdmin).
	Handle(func(c bot.Context, t tele.Context) error {
		menu := createMarkup(c, true)

		t.Send("Please select a user to be removed as an admin.", menu)
		c.SetTextCommand(selectAdminToRemoveEndpoint, nil)

		return nil
	})

var selectAdminToSet = bot.NewCommand(selectAdminToSetEndpoint).
	RestrictUser(bot.RestrictAdmin).
	Handle(func(c bot.Context, t tele.Context) error {
		return setAdminOnSelection(c, t.Message().Text, true)
	})

var selectAdminToRemove = bot.NewCommand(selectAdminToRemoveEndpoint).
	RestrictUser(bot.RestrictAdmin).
	Handle(func(c bot.Context, t tele.Context) error {
		return setAdminOnSelection(c, t.Message().Text, false)
	})

func Attach(router *bot.Router) {

	// setadmin
	router.AddCommand(setAdminCommand)
	router.AddTextCommand(selectAdminToSet)

	// removeadmin
	router.AddCommand(removeAdminCommand)
	router.AddTextCommand(selectAdminToRemove)
}
