package balance

import (
	"fmt"
	"strings"

	"github.com/yardnsm/gohever"
	"github.com/yardnsm/leeches/internal/bot"
	"github.com/yardnsm/leeches/internal/render"

	tele "gopkg.in/telebot.v3"
)

func handleCardBalance(c bot.Context, _ tele.Context, card gohever.CardInterface) error {
	editable, _ := c.SendEditable("⏳ I'm on it...")

	status, err := card.GetStatus()
	if err != nil {
		c.Edit(editable, "I couldn't fetch the card balance.")
		return err
	}

	usageFmt := []string{
		"🛍️ *Keva monthly usage (%d%%):*",
		"",
		"```",
		"%s",
		"```",
		"• 💳 *On card:* %.2f / %d",
		"• 🗓️ *Montly usage:* %.2f / %d",
		"• 💸 *Leftovers:* %.2f",
	}

	if card.Type() == gohever.TypeTeamim {
		usageFmt[0] = "🌮 *Teamim monthly usage (%d%%):*"
	}

	if card.Type() == gohever.TypeSheli {
		usageFmt[0] = "🩲 *Sheli monthly usage (%d%%):*"
	}

	usageViz := render.CardBalance(*status)
	final := fmt.Sprintf(
		strings.Join(usageFmt, "\n"),
		int(100*status.MonthlyUsage),
		usageViz,
		status.CurrentBalance, status.MaxOnCardAmount,
		(float64(status.MaxMonthlyAmount) - status.RemainingMonthlyAmount), status.MaxMonthlyAmount,
		status.Leftovers,
	)

	return c.Edit(editable, final, tele.ModeMarkdown)
}

var balanceKevaCommand = bot.NewCommand(balanceKevaEndpoint).
	Description("View Keva card balance").
	RestrictUser(bot.RestrictApproved).
	Handle(func(c bot.Context, t tele.Context) error {
		return handleCardBalance(c, t, c.Hever.Cards.Keva)
	})

var balanceFoodCommand = bot.NewCommand(balanceFoodEndpoint).
	Description("View Teamim card balance").
	RestrictUser(bot.RestrictApproved).
	Handle(func(c bot.Context, t tele.Context) error {
		return handleCardBalance(c, t, c.Hever.Cards.Teamim)
	})

var balanceSheliCommand = bot.NewCommand(balanceSheliEndpoint).
	Description("View Sheli card balance").
	RestrictUser(bot.RestrictApproved).
	Handle(func(c bot.Context, t tele.Context) error {
		return handleCardBalance(c, t, c.Hever.Cards.Sheli)
	})

var cardTypeToCommand = map[gohever.CardType]*bot.Command{
	gohever.TypeKeva:   balanceKevaCommand,
	gohever.TypeTeamim: balanceFoodCommand,
	gohever.TypeSheli:  balanceSheliCommand,
}

func Attach(router *bot.Router, cards []gohever.CardType) {
	for _, cardType := range cards {
		router.AddCommand(cardTypeToCommand[cardType])
	}
}
