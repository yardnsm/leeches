package commands

import (
	"github.com/yardnsm/leeches/internal/bot"

	"github.com/yardnsm/leeches/internal/commands/admin"
	"github.com/yardnsm/leeches/internal/commands/balance"
	"github.com/yardnsm/leeches/internal/commands/charge"
	"github.com/yardnsm/leeches/internal/commands/users"
)

func Attach(router *bot.Router) {
	admin.Attach(router)
	users.Attach(router)
	balance.Attach(router)
	charge.Attach(router)
}
