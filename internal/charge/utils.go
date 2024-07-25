package charge

import "github.com/yardnsm/leeches/internal/model"

func chargeRequestStateToString(state model.ChargeRequestState) string {
	stateMsg := ""

	switch state {
	case model.StateCreated:
		stateMsg = "❓ Send the request to approval?"
	case model.StateExpired:
		stateMsg = "⏰ Request was expired."
	case model.StatePending:
		stateMsg = "⏳ Request is pending admin approval."
	case model.StateAborted:
		stateMsg = "❌ Request was aborted by the user."
	case model.StateRejected:
		stateMsg = "🛑 Request was rejected by an admin."
	case model.StateApproved:
		stateMsg = "✅ Request was approved by an admin."
	case model.StateCharged:
		stateMsg = "✅ Request was completed successfully."
	case model.StateFailed:
		stateMsg = "🤷 Request was failed for some reason."
	}

	return stateMsg
}
