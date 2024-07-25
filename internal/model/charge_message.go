package model

import (
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

// Represents a stored charge request message
type ChargeMessage struct {
	gorm.Model
	tele.StoredMessage

	UserID uint
	User User

	ChargeRequestID uint
}
