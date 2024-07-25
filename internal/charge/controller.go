package charge

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yardnsm/gohever"
	"github.com/yardnsm/leeches/internal/bot"
	"github.com/yardnsm/leeches/internal/model"

	tele "gopkg.in/telebot.v3"
)

type CallbackEndpoints struct {
	Refresh string
	Reject  string
	Approve string
}

// Manages the process of creating / refetching / updating charge request messages sent to regular
// users and admins
type Controller struct {

	// The bot.Context of some sort. Please note that the given context does not have to be related
	// to the same chat's stored message we're about to edit.
	c *bot.Context

	// The charge request to base the message upon. Please make sure to pass a reference since the
	// controller may update the model.
	request *model.ChargeRequest

	// Import cycles are a bitch.
	callbackEndpoints CallbackEndpoints
}

// Create a new controller for a charge request
func NewController(c *bot.Context, request *model.ChargeRequest, callbackEndpoints CallbackEndpoints) *Controller {
	return &Controller{
		c:                 c,
		request:           request,
		callbackEndpoints: callbackEndpoints,
	}
}

// Get all the users participating this request
func (ctrl *Controller) getParticipants() []model.User {
	requester := ctrl.request.Requester

	// Requester is always a participant
	participants := []model.User{requester}

	// Admins participating only after the state was pending for approval
	if len(ctrl.request.ChargeMessages) <= 1 {
		switch ctrl.request.State {
		case model.StateCreated:
			fallthrough
		case model.StateAborted:
			return participants
		}
	}

	allUsers, _ := ctrl.c.Users.GetAll()
	for _, user := range allUsers {

		// Only take admins that are not the current user
		if user.ID != requester.ID && user.IsAdmin {
			participants = append(participants, user)
		}
	}

	return participants
}

// Get the tele.StoredMessage referece for a participant, if exists
func (ctrl *Controller) getChargeMessageForParticipant(participant model.User) *model.ChargeMessage {
	for _, msg := range ctrl.request.ChargeMessages {
		if msg.User.ID == participant.ID {
			return &msg
		}
	}

	return nil
}

func (ctrl *Controller) addChargeMessage(participant model.User, editable tele.Editable) *model.ChargeMessage {
	storedMessage := bot.EditableToStoredMessage(editable)
	chargeMessage := &model.ChargeMessage{
		User: participant,
	}

	chargeMessage.MessageID = storedMessage.MessageID
	chargeMessage.ChatID = storedMessage.ChatID

	ctrl.request.ChargeMessages = append(ctrl.request.ChargeMessages, *chargeMessage)

	return chargeMessage
}

// Fetch the data from Hever's API and return struct ready for render
func (ctrl *Controller) fetch() (*gohever.CardStatus, *gohever.CardEstimate, error) {
	card := ctrl.c.Hever.Cards.Keva
	if ctrl.request.CardType == gohever.TypeTeamim {
		card = ctrl.c.Hever.Cards.Teamim
	}

	status, err := card.GetStatus()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to fetch card's status: %w", err)
	}

	estimate, err := status.Estimate(float64(ctrl.request.Amount))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to charge with given amount: %w", err)
	}

	return status, estimate, nil
}

// Create the message body for a participant
func (ctrl *Controller) renderBodyForParticipant(
	participant model.User,
	status gohever.CardStatus,
	estimate gohever.CardEstimate,
) string {
	request := ctrl.request
	isAdmin := participant.IsAdmin

	estimateFmt := []string{
		"🟡 *Charge request for Keva:*",
		"",
		"%s requested to charge *%d₪* for %s.",
		"After discount it'll be *%.2f₪*.",
		"",
		"```",
		"%s = %.2f",
		"```",
		"%s",
	}

	estimateFormula := []string{}

	if estimate.Leftovers > 0 {
		estimateFormula = append(
			estimateFormula,
			fmt.Sprintf("%.2f*%.1f", estimate.Leftovers, status.Factors[len(status.Factors)-1].Factor),
		)
	}

	for _, factor := range estimate.Factors {
		if factor.Amount == 0 {
			continue
		}

		estimateFormula = append(estimateFormula, fmt.Sprintf("%.2f*%.2f", factor.Amount, factor.Factor))
	}

	// Update title
	if request.CardType == gohever.TypeTeamim {
		estimateFmt[0] = "🔵 *Charge request for Teamim:*"
	}

	requester := "You"
	stateString := chargeRequestStateToString(request.State)

	if isAdmin {
		if request.RequesterID != participant.ID {
			requester = fmt.Sprintf("*%s*", request.Requester.DisplayName)
		}

		if request.State == model.StateCreated {
			stateString = "🤨 Approve this request?"
		}
	}

	final := fmt.Sprintf(
		strings.Join(estimateFmt, "\n"),
		requester, request.Amount, request.Reason,
		estimate.TotalFactored,
		strings.Join(estimateFormula, " + "), estimate.TotalFactored,
		stateString,
	)

	return final
}

// Create the message markup for a participant
func (ctrl *Controller) renderMarkupForParticipant(participant model.User) *tele.ReplyMarkup {
	state := ctrl.request.State
	if state != model.StateCreated && state != model.StatePending {
		return nil
	}

	selector := &tele.ReplyMarkup{}
	data := strconv.Itoa(int(ctrl.request.ID))

	btnReject := selector.Data("❌", ctrl.callbackEndpoints.Reject, data)
	btnRefresh := selector.Data("🔄", ctrl.callbackEndpoints.Refresh, data)
	btnApprove := selector.Data("✅", ctrl.callbackEndpoints.Approve, data)

	row := selector.Row(btnReject, btnRefresh)

	// Add the approve button only for admins or users when the request is created
	if participant.IsAdmin || ctrl.request.State == model.StateCreated {
		row = append(row, btnApprove)
	}

	selector.Inline(
		row,
	)

	return selector
}

// Save (upsert) the model
func (ctrl *Controller) Save() error {
	return ctrl.c.ChargeRequests.Save(ctrl.request)
}

// Update the charge request for all the participating users based on the model only
func (ctrl *Controller) render() error {
	participants := ctrl.getParticipants()

	status := ctrl.request.CachedCardStatus
	estimate := ctrl.request.CachedCardEstimate

	for _, participant := range participants {
		body := ctrl.renderBodyForParticipant(participant, status, estimate)
		markup := ctrl.renderMarkupForParticipant(participant)

		// Get the stored message for the user if exists
		msg := ctrl.getChargeMessageForParticipant(participant)

		if msg != nil {

			// Update existing
			err := ctrl.c.Edit(msg, body, tele.ModeMarkdown, markup)
			if err != nil && !errors.Is(err, tele.ErrSameMessageContent) {
				return err
			}
		} else {

			// Send a new message to the participant
			msg, err := ctrl.c.Bot().Send(&participant, body, tele.ModeMarkdown, markup)
			if err != nil {
				return err
			}

			// Store the message in the model
			ctrl.addChargeMessage(participant, msg)
		}
	}

	// Save the request, we may have updated the charge messages
	return ctrl.Save()
}

// Initialize the charge request: create the request in the database and send the message to all the
// participating users.
func (ctrl *Controller) Init(editable tele.Editable) error {
	// Create and store the message for the requester
	ctrl.addChargeMessage(ctrl.request.Requester, editable)

	// Perform a hard update: fetch from the API and render
	return ctrl.HardUpdate()
}

// Re-fetches the data from Hever's API, update the model but does not cause a re-render. This
// function may return an error from the estimation.
// Will also re-save the model.
func (ctrl *Controller) Refetch() error {
	status, estimate, err := ctrl.fetch()
	if err != nil {
		return err
	}

	ctrl.request.CachedCardStatus = *status
	ctrl.request.CachedCardEstimate = *estimate
	ctrl.request.CachedAt = time.Now()

	return ctrl.Save()
}

// Same as render()
func (ctrl *Controller) SoftUpdate() error {
	return ctrl.render()
}

// Update the charge request for all the participating users and re-fetch the data from the Hever's
// API. Will also re-save the model.
func (ctrl *Controller) HardUpdate() error {
	err := ctrl.Refetch()
	if err != nil {
		return err
	}

	return ctrl.render()
}
