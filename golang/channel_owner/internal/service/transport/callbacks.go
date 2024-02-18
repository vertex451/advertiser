package transport

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type CallBackQueryParams struct {
	ChatID   int64
	Page     string
	Variable string
}

func (t *Transport) handleCallbackQuery(query *tgbotapi.CallbackQuery) *tgbotapi.MessageConfig {
	params := parseCallBackQuery(query)

	if params.Page != Back {
		t.addCrumbs(params)
	}

	return t.NavigateToPage(params)
}

func (t *Transport) NavigateToPage(params CallBackQueryParams) *tgbotapi.MessageConfig {
	switch params.Page {
	case Start:
		return t.start(params.ChatID)
	case Back:
		return t.back(params.ChatID)

	case AllTopics:
		return t.allTopics(params.ChatID)
	case MyChannels:
		return t.listMyChannels(params.ChatID)
	case ListChannelsTopics:
		return t.listChannelTopics(params.ChatID, params.Variable)
	case EditChannelsTopics:
		return t.editTopicsPrompt(params.ChatID, params.Variable)

	default:
		return t.start(params.ChatID)
	}
}

func parseCallBackQuery(query *tgbotapi.CallbackQuery) CallBackQueryParams {
	parsed := strings.Split(query.Data, "/")

	res := CallBackQueryParams{
		ChatID: query.Message.Chat.ID,
		Page:   parsed[0],
	}

	if len(parsed) > 1 {
		res.Variable = parsed[1]
	}

	return res
}
