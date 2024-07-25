package users

import (
	"fmt"
	"strconv"

	"github.com/yardnsm/leeches/internal/bot"
	"github.com/yardnsm/leeches/internal/model"

	tele "gopkg.in/telebot.v3"
)

var addUserCommand = bot.NewCommand(addUserEndpoint).
	Description("Add a user to the bot").
	RestrictUser(bot.RestrictAdmin).
	Handle(func(c bot.Context, t tele.Context) error {
		var state addUserState

		t.Send("What's the user id?")
		c.SetTextCommand(getUserIdEndpoint, state)

		return nil
	})

var removeUserCommand = bot.NewCommand(removeUserEndpoint).
	Description("Remove a user from the bot").
	RestrictUser(bot.RestrictAdmin).
	Handle(func(c bot.Context, t tele.Context) error {
		menu := createUsersMarkup(c)

		t.Send("Please select a user to remove", menu)
		c.SetTextCommand(selectUserToRemoveEndpoint, nil)

		return nil
	})

var getUserId = bot.NewCommand(getUserIdEndpoint).
	RestrictUser(bot.RestrictAdmin).
	Handle(bot.CreateStatefulHandler(
		func(c bot.Context, t tele.Context, state addUserState) (interface{}, *addUserState, error) {
			state.userID = t.Message().Text
			t.Send("What is the display name for the user?")
			return getDisplayNameEndpoint, &state, nil
		},
	))

var getDisplayName = bot.NewCommand(getDisplayNameEndpoint).
	RestrictUser(bot.RestrictAdmin).
	Handle(bot.CreateStatefulHandler(
		func(c bot.Context, t tele.Context, state addUserState) (interface{}, *addUserState, error) {
			state.displayName = t.Message().Text

			telegramID, err := strconv.Atoi(state.userID)
			if err != nil {
				t.Send("An error occured while creating the user. Details:")
				return nil, nil, err
			}

			user := model.User{
				TelegramID:  int64(telegramID),
				DisplayName: state.displayName,
				IsApproved:  true,
				IsAdmin:     false,
			}

			err = c.Users.Create(&user)
			if err != nil {
				t.Send("An error occured while creating the user. Details:")
				return nil, nil, err
			}

			err = t.Send(fmt.Sprintf("New user created! Internal ID: %d", user.ID))
			return nil, nil, err
		},
	))

var selectUserToRemove = bot.NewCommand(selectUserToRemoveEndpoint).
	RestrictUser(bot.RestrictAdmin).
	Handle(func(c bot.Context, t tele.Context) error {
		defer c.ClearTextCommand()

		selection := t.Message().Text
		extracted := selectionRegexp.FindStringSubmatch(selection)

		if len(extracted) == 0 {
			return c.SendAndCloseKeyboard("An invalid selection was made.")
		}

		userID, _ := strconv.ParseUint(extracted[1], 10, 0)

		// Do not allow the current user remove himself.
		if uint(userID) == c.CurrentUser.ID {
			return c.SendAndCloseKeyboard("Suicide is discouraged.")
		}

		err := c.Users.Delete(uint(userID))
		if err != nil {
			c.SendAndCloseKeyboard("An error occured while removing the user. Details:")
			return err
		}

		return c.SendAndCloseKeyboard(fmt.Sprintf("User %s was removed.", selection))
	})

func Attach(router *bot.Router) {

	// adduser
	router.AddCommand(addUserCommand)
	router.AddTextCommand(getUserId)
	router.AddTextCommand(getDisplayName)

	// removeuser
	router.AddCommand(removeUserCommand)
	router.AddTextCommand(selectUserToRemove)
}
