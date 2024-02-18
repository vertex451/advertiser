package transport

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (t *Transport) start(respondTo int64) *tgbotapi.MessageConfig {
	t.resetState(respondTo)

	msg := addNavigationButtons(
		tgbotapi.NewMessage(respondTo, "Choose action:"),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("List my campaigns", fmt.Sprintf("%s", MyCampaigns)),
			tgbotapi.NewInlineKeyboardButtonData("Create new campaign", CreateCampaign),
			tgbotapi.NewInlineKeyboardButtonData("List all topics", fmt.Sprintf("%s", AllTopicsWithCoverage)),
		),
	)
	return &msg
}

func (t *Transport) back(respondTo int64) *tgbotapi.MessageConfig {
	state := t.getState(respondTo)
	if len(state.crumbs) <= 1 {
		return t.start(respondTo)
	}

	params := state.crumbs[len(state.crumbs)-2]
	state.crumbs = state.crumbs[:len(state.crumbs)-1]

	t.setState(respondTo, state)

	return t.NavigateToPage(params)
}

func (t *Transport) allTopics(respondTo int64) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf(`
Supported topics:
%s
`, strings.Join(t.uc.AllTopics(), ", ")))
	msg = addNavigationButtons(msg, nil)

	return &msg
}
