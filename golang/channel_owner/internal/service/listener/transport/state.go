package transport

import (
	"advertiser/shared/pkg/service/transport"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type stateData struct {
	crumbs          []transport.CallBackQueryParams
	state           BotState
	channelID       int64
	botDirectChatID int64
	campaignID      string
	adID            string
}

func (s *Transport) handleStateQuery(update tgbotapi.Update) *tgbotapi.MessageConfig {
	chatID := update.Message.Chat.ID
	state := s.getState(chatID)

	switch state.state {
	case StateEditTopics:
		topics := strings.Split(update.Message.Text, ",")
		return s.editChannelTopics(chatID, state.channelID, topics)
	default:
		return s.start(chatID)
	}
}

// setState sets the state of the conversation for a given chat ID
func (s *Transport) setState(chatID int64, data stateData) {
	s.state.Store(chatID, data)
}

// getState retrieves the state of the conversation for a given chat ID
func (s *Transport) getState(chatID int64) stateData {
	state, ok := s.state.Load(chatID)
	if !ok {
		return stateData{
			state: StateStart,
		}
	}

	return state.(stateData)
}

func (s *Transport) resetState(chatID int64) {
	s.state.Store(chatID, stateData{
		state: StateStart,
	})
}

func (s *Transport) addCrumbs(params transport.CallBackQueryParams) {
	rawState, ok := s.state.Load(params.ChatID)
	var state stateData
	if ok {
		state = rawState.(stateData)
	} else {
		state = stateData{
			state: StateStart,
		}
	}

	state.crumbs = append(state.crumbs, params)

	s.state.Store(params.ChatID, state)
}
