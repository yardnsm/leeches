package admin

import "regexp"

const (
	setAdminEndpoint         = "/setadmin"
	selectAdminToSetEndpoint = "/setadmin/select"

	removeAdminEndpoint         = "/removeadmin"
	selectAdminToRemoveEndpoint = "/removeadmin/select"
)

var (
	selectionRegexp = regexp.MustCompile("\\((\\d+)\\)")
)
