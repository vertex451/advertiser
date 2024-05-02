package transport

import (
	"advertiser/shared/pkg/service/transport"
)

type stateData struct {
	lastMsgID               int
	crumbs                  []transport.CallBackQueryParams
	state                   BotState
	storeInitialChannelData bool
	adChanID                string
	channelID               int64
	campaignID              string
	adID                    string
}

// setState sets the state of the conversation for a given chat ID
func (s *Service) setState(chatID int64, data stateData) {
	s.state.Store(chatID, data)
}

// getState retrieves the state of the conversation for a given chat ID
func (s *Service) getState(chatID int64) stateData {
	state, ok := s.state.Load(chatID)
	if !ok {
		return stateData{
			state: StateStart,
		}
	}

	return state.(stateData)
}

func (s *Service) resetState(chatID int64) {
	state := s.getState(chatID)
	s.state.Store(chatID, stateData{
		state:     StateStart,
		lastMsgID: state.lastMsgID,
	})
}

func (s *Service) addCrumbs(params transport.CallBackQueryParams) {
	rawState, ok := s.state.Load(params.UserID)
	var state stateData
	if ok {
		state = rawState.(stateData)
	} else {
		state = stateData{
			state: StateStart,
		}
	}

	state.crumbs = append(state.crumbs, params)

	s.state.Store(params.UserID, state)
}
