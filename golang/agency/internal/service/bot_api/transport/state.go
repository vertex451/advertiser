package transport

import (
	"advertiser/shared/pkg/service/transport"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type stateData struct {
	lastMsgID       int
	crumbs          []transport.CallBackQueryParams
	state           BotState
	channelID       int64
	botDirectChatID int64
	campaignID      string
	adID            string
}

func (t *Transport) handleStateQuery(update tgbotapi.Update) *transport.Msg {
	chatID := update.Message.Chat.ID
	state := t.getState(chatID)

	switch state.state {
	case StateSetCampaignName:
		return t.createCampaign(chatID, update.Message.Text)
	case StateCreateAd:
		return t.upsertAd(chatID, state.campaignID, "", update.Message.Text)
	case StateUpdateAd:
		return t.upsertAd(chatID, "", state.adID, update.Message.Text)
	default:
		return t.start(chatID)
	}
}

// setState sets the state of the conversation for a given chat ID
func (t *Transport) setState(chatID int64, data stateData) {
	t.state.Store(chatID, data)
}

// getState retrieves the state of the conversation for a given chat ID
func (t *Transport) getState(chatID int64) stateData {
	state, ok := t.state.Load(chatID)
	if !ok {
		return stateData{
			state: StateStart,
		}
	}

	return state.(stateData)
}

func (t *Transport) resetState(chatID int64) {
	state := t.getState(chatID)
	t.state.Store(chatID, stateData{
		state:     StateStart,
		lastMsgID: state.lastMsgID,
	})
}

func (t *Transport) addCrumbs(params transport.CallBackQueryParams) {
	rawState, ok := t.state.Load(params.UserID)
	var state stateData
	if ok {
		state = rawState.(stateData)
	} else {
		state = stateData{
			state: StateStart,
		}
	}

	state.crumbs = append(state.crumbs, params)

	t.state.Store(params.UserID, state)
}
