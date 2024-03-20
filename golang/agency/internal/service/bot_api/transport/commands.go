package transport

import (
	"advertiser/shared/pkg/service/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (t *Transport) start(respondTo int64) types.CustomMessage {
	t.resetState(respondTo)

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, "Choose action:"),
		[][]tgbotapi.InlineKeyboardButton{{
			tgbotapi.NewInlineKeyboardButtonData("My campaigns", fmt.Sprintf("%s", MyCampaigns)),
			tgbotapi.NewInlineKeyboardButtonData("Create a new campaign", CreateCampaign),
			tgbotapi.NewInlineKeyboardButtonData("All topics", fmt.Sprintf("%s", AllTopicsWithCoverage)),
		}},
		true,
		false,
	)
}

func (t *Transport) back(respondTo int64) types.CustomMessage {
	state := t.getState(respondTo)
	if len(state.crumbs) <= 1 {
		return t.start(respondTo)
	}

	params := state.crumbs[len(state.crumbs)-2]
	state.crumbs = state.crumbs[:len(state.crumbs)-1]
	state.state = StateStart

	t.setState(respondTo, state)

	return t.NavigateToPage(params)
}

func (t *Transport) allTopics(respondTo int64) types.CustomMessage {
	msgText := fmt.Sprintf(`
Supported topics:
%s
`, strings.Join(t.uc.AllTopics(), ", "))

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, msgText),
		nil,
		true,
		false,
	)
}
