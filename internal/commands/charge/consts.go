package charge

import (
	"errors"
	"time"
)

const (
	chargeKevaEndpoint = "/chargekeva"
	chargeFoodEndpoint = "/chargefood"

	chargeCardGetAmountEndpoint = "/chargecard/getamount"
	chargeCardGetReasonEndpoint = "/chargecard/getreason"

	// For some reason the callback cannot contain any slashes...
	chargeCardRefreshEndpoint = "chargecardrefresh"
	chargeCardRejectEndpoint  = "chargecardreject"
	chargeCardApproveEndpoint = "chargecardapprove"
)

// Errors
var (
	ErrRequestNotFound  = errors.New("unable to find the charge request")
	ErrUnauthorizedUser = errors.New("unauthorized attempt at performing the action")

	ErrNotInValidState = errors.New("charge request should be in pending state before approving")
	ErrExpired         = errors.New("charge request expired")
)
