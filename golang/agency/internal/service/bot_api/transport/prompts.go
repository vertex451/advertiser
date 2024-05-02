package transport

import (
	"advertiser/shared/pkg/service/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Transport) createCampaignPrompt(respondTo int64) types.CustomMessage {
	msg := tgbotapi.NewMessage(respondTo, "Enter campaign name:")

	t.setState(respondTo, stateData{
		state: StateSetCampaignName,
	})

	return types.NewCustomMessageConfig(
		msg,
		nil,
		false,
		true,
	)
}

func (t *Transport) createAdPrompt(respondTo int64, variable string, state BotState) types.CustomMessage {
	var action string
	switch state {
	case StateCreateAd:
		action = "create"
		t.setState(respondTo, stateData{
			state:      StateCreateAd,
			campaignID: variable,
		})
	}

	msgText := fmt.Sprintf(`
To %s an advertisement, send a message in the following format:
Name: McDonald's
TargetTopics: food
BudgetUSD: 100
CostPerMile: 0.1
`, action)

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, msgText),
		nil,
		false,
		true,
	)
}
