package transport

import (
	"advertiser/shared/pkg/service/transport"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Transport) createCampaignPrompt(respondTo int64) *transport.Msg {
	msg := tgbotapi.NewMessage(respondTo, "Send new campaign name:")

	t.setState(respondTo, stateData{
		state:           StateSetCampaignName,
		botDirectChatID: respondTo,
	})

	return &transport.Msg{
		Msg: msg,
	}
}

func (t *Transport) upsertAdPrompt(respondTo int64, variable string, state BotState) *transport.Msg {
	var msg tgbotapi.MessageConfig
	var action string
	switch state {
	case StateCreateAd:
		action = "create"
		t.setState(respondTo, stateData{
			state:      StateCreateAd,
			campaignID: variable,
		})
	case StateUpdateAd:
		action = "update"
		t.setState(respondTo, stateData{
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
	msg = tgbotapi.NewMessage(respondTo, promptText)

	return &transport.Msg{
		Msg: msg,
	}
}
