package users

import "regexp"

const (
	addUserEndpoint        = "/adduser"
	getUserIdEndpoint      = "/adduser/getuserid"
	getDisplayNameEndpoint = "/adduser/getdisplayname"

	removeUserEndpoint         = "/removeuser"
	selectUserToRemoveEndpoint = "/removeuser/select"
)

var (
	selectionRegexp = regexp.MustCompile("\\((\\d+)\\)")
)

type addUserState struct {
	userID      string
	displayName string
}
