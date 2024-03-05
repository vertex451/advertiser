package transport

import (
	"advertiser/shared/pkg/service/transport"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (t *Transport) start(respondTo int64) *transport.Msg {
	t.resetState(respondTo)

	msg := transport.AddNavigationButtons(
		tgbotapi.NewMessage(respondTo, "Choose action:"),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("My campaigns", fmt.Sprintf("%s", MyCampaigns)),
			tgbotapi.NewInlineKeyboardButtonData("Create a new campaign", CreateCampaign),
			tgbotapi.NewInlineKeyboardButtonData("All topics", fmt.Sprintf("%s", AllTopicsWithCoverage)),
		),
	)

	return &transport.Msg{
		Msg: msg,
	}
}

func (t *Transport) back(respondTo int64) *transport.Msg {
	state := t.getState(respondTo)
	if len(state.crumbs) <= 1 {
		return t.start(respondTo)
	}

	params := state.crumbs[len(state.crumbs)-2]
	state.crumbs = state.crumbs[:len(state.crumbs)-1]

	t.setState(respondTo, state)

	return t.NavigateToPage(params)
}

func (t *Transport) allTopics(respondTo int64) *transport.Msg {
	var msg tgbotapi.MessageConfig
	msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf(`
Supported topics:
%s
`, strings.Join(t.uc.AllTopics(), ", ")))
	msg = transport.AddNavigationButtons(msg, nil)

	return &transport.Msg{
		Msg: msg,
	}
}
