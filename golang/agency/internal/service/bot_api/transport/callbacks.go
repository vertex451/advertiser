package transport

import (
	"advertiser/shared/pkg/service/transport"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Transport) handleCallbackQuery(query *tgbotapi.CallbackQuery) *transport.Msg {
	params := transport.ParseCallBackQuery(query)

	if params.Page != Back {
		t.addCrumbs(params)
	}

	return t.NavigateToPage(params)
}

func (t *Transport) NavigateToPage(params transport.CallBackQueryParams) *transport.Msg {
	switch params.Page {
	case Start:
		return t.start(params.ChatID)
	case Back:
		return t.back(params.ChatID)

	case AllTopicsWithCoverage:
		return t.allTopicsWithCoverage(params.ChatID)
	case CreateCampaign:
		return t.createCampaignPrompt(params.ChatID)
	case MyCampaigns:
		return t.listMyCampaigns(params.ChatID)
	case CampaignDetails:
		return t.campaignDetails(params.ChatID, params.Variable)
	case CreateAd:
		return t.upsertAdPrompt(params.ChatID, params.Variable, StateCreateAd)
	case AdDetails:
		return t.GetAdDetails(params.ChatID, params.Variable)
	case EditAd:
		return t.upsertAdPrompt(params.ChatID, params.Variable, StateUpdateAd)
	case RunAd:
		return t.RunAd(params.ChatID, params.Variable)

	default:
		return t.start(params.ChatID)
	}
}
