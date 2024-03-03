package transport

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/transport"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Transport) handleCallbackQuery(query *tgbotapi.CallbackQuery) *transport.Msg {
	params := transport.ParseCallBackQuery(query)

	if params.Page != constants.Back {
		s.addCrumbs(params)
	}

	return s.NavigateToPage(params)
}

func (s *Transport) NavigateToPage(params transport.CallBackQueryParams) *transport.Msg {
	switch params.Page {
	case constants.Start:
		return s.start(params.ChatID)
	case constants.Back:
		return s.back(params.ChatID)

	case constants.AllTopics:
		return s.allTopics(params.ChatID)
	case Moderate:
		return s.moderate(params.ChatID)
	case MyChannels:
		return s.listMyChannels(params.ChatID)
	case ListChannelsTopics:
		return s.listChannelTopics(params.ChatID, params.Variable)
	case EditChannelsTopics:
		return s.editTopicsPrompt(params.ChatID, params.Variable)
	case ModerateDetails:
		return s.GetAdvertisementDetails(params.ChatID, params.Variable)
	case PostNow:
		return s.moderationDecision(params.ChatID, PostNow, params.Variable)
	case RejectAd:
		return s.moderationDecision(params.ChatID, RejectAd, params.Variable)

	default:
		return s.start(params.ChatID)
	}
}
