package transport

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/transport"
	"advertiser/shared/pkg/service/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Transport) handleCallbackQuery(query *tgbotapi.CallbackQuery) types.CustomMessage {
	params := transport.ParseCallBackQuery(query)

	if params.Page != constants.Back {
		t.addCrumbs(params)
	}

	return t.NavigateToPage(params)
}

func (t *Transport) NavigateToPage(params transport.CallBackQueryParams) types.CustomMessage {
	switch params.Page {
	case constants.Start:
		return t.start(params.UserID)
	case constants.Back:
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
		return t.createAdPrompt(params.UserID, params.Variable, StateCreateAd)
	case AdDetails:
		return t.getAdDetails(params.UserID, params.Variable)
	case constants.ViewAdMessage:
		return t.viewAdMessage(params.UserID, params.Variable)
	//case EditAd:
	//	return t.createAdPrompt(params.UserID, params.Variable, StateUpdateAd)
	case RunAd:
		return t.RunAd(params.UserID, params.Variable)

	default:
		return t.start(params.UserID)
	}
}
