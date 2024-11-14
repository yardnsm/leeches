package charge

import (
	"strconv"

	"github.com/yardnsm/gohever"
	"github.com/yardnsm/leeches/internal/bot"
	"github.com/yardnsm/leeches/internal/charge"
	"github.com/yardnsm/leeches/internal/model"

	tele "gopkg.in/telebot.v3"
)

// This state represents the process of creating a charge request
type chargeCardState struct {
	cardType gohever.CardType
	amount   int32
	reason   string
}

func createController(c *bot.Context, request *model.ChargeRequest) *charge.Controller {
	return charge.NewController(
		c,
		request,
		charge.CallbackEndpoints{
			Refresh: chargeCardRefreshEndpoint,
			Reject:  chargeCardRejectEndpoint,
			Approve: chargeCardApproveEndpoint,
		},
	)
}

func handleChargeCardStart(c bot.Context, t tele.Context, cardType gohever.CardType) error {
	state := chargeCardState{
		cardType: cardType,
	}

	t.Send("How much would you like to charge?", tele.ForceReply)
	c.SetTextCommand(chargeCardGetAmountEndpoint, state)

	return nil
}

func handleChargeCardCreate(c bot.Context, t tele.Context, state chargeCardState) error {
	// Create a message, which will be used as a notifier for the user
	editable, err := c.SendEditable("🧑‍🏫 I'm calculating things...")
	if err != nil {
		return err
	}

	// Create a base charge request
	request := model.ChargeRequest{
		Amount: state.amount,
		Reason: state.reason,

		CardType: state.cardType,
		State:    model.StateCreated,

		Requester: *c.CurrentUser,
	}

	ctrl := createController(&c, &request)

	err = ctrl.Init(editable)
	if err != nil {
		c.Edit(editable, "I couldn't create a charge request.")
		return err
	}

	return nil
}

var chargeKevaCommand = bot.NewCommand(chargeKevaEndpoint).
	Description("Charge the keva card").
	RestrictUser(bot.RestrictApproved).
	Handle(func(c bot.Context, t tele.Context) error {
		return handleChargeCardStart(c, t, gohever.TypeKeva)
	})

var chargeFoodCommand = bot.NewCommand(chargeFoodEndpoint).
	Description("Charge the Teamim card").
	RestrictUser(bot.RestrictApproved).
	Handle(func(c bot.Context, t tele.Context) error {
		return handleChargeCardStart(c, t, gohever.TypeTeamim)
	})

var chargeSheliCommand = bot.NewCommand(chargeSheliEndpoint).
	Description("Charge the Sheli card").
	RestrictUser(bot.RestrictApproved).
	Handle(func(c bot.Context, t tele.Context) error {
		return handleChargeCardStart(c, t, gohever.TypeSheli)
	})

var chargeCardGetAmount = bot.NewCommand(chargeCardGetAmountEndpoint).
	RestrictUser(bot.RestrictApproved).
	Handle(bot.CreateStatefulHandler(
		func(c bot.Context, t tele.Context, state chargeCardState) (interface{}, *chargeCardState, error) {
			amount, err := strconv.Atoi(t.Message().Text)
			if err != nil {
				t.Send("Please type a valid amount.")
				return chargeCardGetAmountEndpoint, &state, nil
			}

			state.amount = int32(amount)

			t.Send("Why you want to charge the card?")
			return chargeCardGetReasonEndpoint, &state, nil
		},
	))

var chargeCardGetReason = bot.NewCommand(chargeCardGetReasonEndpoint).
	RestrictUser(bot.RestrictApproved).
	Handle(bot.CreateStatefulHandler(
		func(c bot.Context, t tele.Context, state chargeCardState) (interface{}, *chargeCardState, error) {
			state.reason = t.Message().Text
			return nil, nil, handleChargeCardCreate(c, t, state)
		},
	))

var chargeCardRefreshCallback = bot.NewCommand(chargeCardRefreshEndpoint).
	RestrictUser(bot.RestrictApproved).
	Handle(func(c bot.Context, t tele.Context) error {
		return handleRefreshCallback(c, t)
	})

var chargeCardRejectCallback = bot.NewCommand(chargeCardRejectEndpoint).
	RestrictUser(bot.RestrictApproved).
	Handle(func(c bot.Context, t tele.Context) error {
		return handleRejectCallback(c, t)
	})

var chargeCardApproveCallback = bot.NewCommand(chargeCardApproveEndpoint).
	RestrictUser(bot.RestrictApproved).
	Handle(func(c bot.Context, t tele.Context) error {
		return handleApproveCalback(c, t)
	})

var cardTypeToCommand = map[gohever.CardType]*bot.Command{
	gohever.TypeKeva:   chargeKevaCommand,
	gohever.TypeTeamim: chargeFoodCommand,
	gohever.TypeSheli:  chargeSheliCommand,
}

func Attach(router *bot.Router, cards []gohever.CardType) {
	for _, cardType := range cards {
		router.AddCommand(cardTypeToCommand[cardType])
	}

	router.AddCallback(chargeCardRefreshCallback)
	router.AddCallback(chargeCardRejectCallback)
	router.AddCallback(chargeCardApproveCallback)

	router.AddTextCommand(chargeCardGetAmount)
	router.AddTextCommand(chargeCardGetReason)
}
