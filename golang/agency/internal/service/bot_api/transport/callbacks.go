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
		return t.start(params.UserID)
	case Back:
		return t.back(params.UserID)

	case AllTopicsWithCoverage:
		return t.allTopicsWithCoverage(params.UserID)
	case CreateCampaign:
		return t.createCampaignPrompt(params.UserID)
	case MyCampaigns:
		return t.listMyCampaigns(params.UserID)
	case CampaignDetails:
		return t.campaignDetails(params.UserID, params.Variable)
	case CreateAd:
		return t.upsertAdPrompt(params.UserID, params.Variable, StateCreateAd)
	case AdDetails:
		return t.GetAdDetails(params.UserID, params.Variable)
	case EditAd:
		return t.upsertAdPrompt(params.UserID, params.Variable, StateUpdateAd)
	case RunAd:
		return t.RunAd(params.UserID, params.Variable)

	default:
		return t.start(params.UserID)
	}
}
