package bot

import (
	"errors"

	"github.com/yardnsm/leeches/internal/model"
)

type restrictFunc func(*model.User) error

func RestrictAdmin(user *model.User) error {
	if user != nil && user.IsAdmin {
		return nil
	}

	return errors.New("")
}

func RestrictApproved(user *model.User) error {
	if user != nil && user.IsApproved {
		return nil
	}

	return errors.New("")
}

func RestrictNone(user *model.User) error {
	return nil
}
