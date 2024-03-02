package transport

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Transport) createCampaignPrompt(responseTo int64) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(responseTo, "Send new campaign name:")

	t.setState(responseTo, stateData{
		state:           StateSetCampaignName,
		botDirectChatID: responseTo,
	})

	return &msg
}

func (t *Transport) upsertAdPrompt(responseTo int64, variable string, state BotState) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	var action string
	switch state {
	case StateCreateAd:
		action = "create"
		t.setState(responseTo, stateData{
			state:      StateCreateAd,
			campaignID: variable,
		})
	case StateUpdateAd:
		action = "update"
		t.setState(responseTo, stateData{
			state: StateUpdateAd,
			adID:  variable,
		})
	}

	promptText := fmt.Sprintf(`
To %s an advertisement, send a message in the following format:
Name: MyAwesomeAdvertisement
TargetTopics: topic1, topic2, topic3
BudgetUSD: 100
CostPerView: 0.1
Message: Follow this [link](https://www.investing.com/) to find more about investments!
`, action)
	msg = tgbotapi.NewMessage(responseTo, promptText)

	return &msg
}
