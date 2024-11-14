package commands

import (
	"github.com/yardnsm/gohever"
	"github.com/yardnsm/leeches/internal/bot"

	"github.com/yardnsm/leeches/internal/commands/admin"
	"github.com/yardnsm/leeches/internal/commands/balance"
	"github.com/yardnsm/leeches/internal/commands/charge"
	"github.com/yardnsm/leeches/internal/commands/users"
)

func Attach(router *bot.Router, cards []gohever.CardType) {
	admin.Attach(router)
	users.Attach(router)

	balance.Attach(router, cards)
	charge.Attach(router, cards)
}
