package charge

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/yardnsm/gohever"
	"github.com/yardnsm/leeches/internal/bot"
	"github.com/yardnsm/leeches/internal/charge"
	"github.com/yardnsm/leeches/internal/config"
	"github.com/yardnsm/leeches/internal/model"

	tele "gopkg.in/telebot.v3"
)

var mu sync.Mutex

// This function is called at the begining of each callback
func extractCallbackParams(c bot.Context, t tele.Context) (
	*model.ChargeRequest,
	*charge.Controller,
	bool, // isAdmin
	bool, // isCurrentUserRequester
	error,
) {
	id, _ := strconv.Atoi(t.Callback().Data)

	request, err := c.ChargeRequests.GetByID(uint(id))
	if err != nil {
		return nil, nil, false, false, err
	}

	if request == nil {
		return nil, nil, false, false, ErrRequestNotFound
	}

	isAdmin := c.CurrentUser.IsAdmin
	isCurrentUserRequester := request.Requester.ID == c.CurrentUser.ID

	// Cases when the callback is "somehow" called from a different user
	if !isAdmin && !isCurrentUserRequester {
		return nil, nil, false, false, ErrUnauthorizedUser
	}

	ctrl := createController(&c, request)

	return request, ctrl, isAdmin, isCurrentUserRequester, nil
}

// Verify a charge request. Does not update the model.
func verifyChargeRequest(request model.ChargeRequest) error {
	if request.State != model.StatePending {
		return ErrNotInValidState
	}

	if time.Now().After(request.CreatedAt.Add(config.MaxChargeRequestTime)) {
		return ErrExpired
	}

	return nil
}

func handleRefreshCallback(c bot.Context, t tele.Context) error {
	mu.Lock()
	defer mu.Unlock()

	request, ctrl, _, _, err := extractCallbackParams(c, t)
	if err != nil {
		return err
	}

	// Verify request expiry
	err = verifyChargeRequest(*request)
	if errors.Is(err, ErrExpired) {
		request.State = model.StateExpired
	}

	err = ctrl.HardUpdate()
	if err != nil {
		t.Respond(&tele.CallbackResponse{
			Text: "Something went wrong when refreshing",
		})

		return err
	}

	return t.Respond(&tele.CallbackResponse{
		Text: "Refreshed successfully",
	})
}

func handleRejectCallback(c bot.Context, t tele.Context) error {
	mu.Lock()
	defer mu.Unlock()

	request, ctrl, isAdmin, isCurrentUserRequester, err := extractCallbackParams(c, t)
	if err != nil {
		return err
	}

	// An abort callback was called. We'll abort it and set the state according to the user
	// aborted it.
	request.State = model.StateAborted

	// Admin can reject all and the status is changed to rejected
	if isAdmin && !isCurrentUserRequester {
		request.State = model.StateRejected
	}

	err = ctrl.SoftUpdate()
	if err != nil {
		return err
	}

	return t.Respond()
}

func handleApproveCalback(c bot.Context, t tele.Context) error {
	mu.Lock()
	defer mu.Unlock()

	request, ctrl, isAdmin, _, err := extractCallbackParams(c, t)
	if err != nil {
		return err
	}

	// Save and re-render after finishing
	defer (func() error {
		err := ctrl.Save()
		if err != nil {
			return err
		}

		return ctrl.SoftUpdate()
	})()

	if isAdmin {

		// Make the request pending at first
		request.State = model.StatePending

		// We'll verify first since it's lighter than the rest
		err = verifyChargeRequest(*request)
		switch err {
		case ErrNotInValidState:
			return t.Respond(&tele.CallbackResponse{Text: "You cannot approve charge requests with state other than pending."})
		case ErrExpired:
			request.State = model.StateExpired
			return t.Respond(&tele.CallbackResponse{Text: "The charge request has expired."})
		}

		prevTotalFactored := request.CachedCardEstimate.TotalFactored

		// We'll re-fetch the request in order to get a fresh data
		err = ctrl.Refetch()
		if err != nil {
			return err
		}

		newTotalFactored := request.CachedCardEstimate.TotalFactored
		if newTotalFactored != prevTotalFactored {
			return t.Respond(&tele.CallbackResponse{Text: "The estimation has changed - please approve again."})
		}

		card := c.GetCardByType(request.CardType)

		// Perform the charge!
		result, err := card.Load(request.CachedCardStatus, request.Amount)
		if err != nil {
			request.State = model.StateFailed
			t.Respond(&tele.CallbackResponse{Text: "Charge failed."})

			return err
		}

		if result.Status == gohever.StatusSuccess {
			request.State = model.StateCharged
			return t.Respond(&tele.CallbackResponse{Text: result.RawMessage})
		}

		if result.Status == gohever.StatusError {
			request.State = model.StateFailed
			t.Send(result.RawMessage)
			return t.Respond(&tele.CallbackResponse{Text: "Charge failed."})
		}

		return nil;
	}

	// Update state. This will allow the controller to send charge mesages to admins as well
	request.State = model.StatePending

	return t.Respond(&tele.CallbackResponse{Text: "Request sent to admins"})
}
