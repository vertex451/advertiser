package integration_tests

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/repo/models"
	mocks "advertiser/shared/tg_bot_api/mocks"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"tg-bot/internal/service/bot_api/transport"
)

type channel struct {
	id int64
}

type bot struct {
	id int64
}

type testCase struct {
	testName        string
	preHook         func(db *gorm.DB) string
	update          func(string) *tgbotapi.Update
	expectedMsgText string
	expectedButtons []tgbotapi.InlineKeyboardButton
}

var (
	testChannel = channel{id: -1002093237940}
	testBot     = bot{id: 6406834985}
)

func startCommandUpdate(data string) *tgbotapi.Update {
	return &tgbotapi.Update{
		Message: &tgbotapi.Message{
			From:     &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: len(constants.Start) + 1}},
			Text:     fmt.Sprintf("/%s", constants.Start)},
	}
}

func startCallbackUpdate(data string) *tgbotapi.Update {
	return &tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Data: constants.Start,
		},
	}
}

func allTopicsCommandUpdate(data string) *tgbotapi.Update {
	return &tgbotapi.Update{
		Message: &tgbotapi.Message{
			From:     &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: len(constants.AllTopics) + 1}},
			Text:     fmt.Sprintf("/%s", constants.AllTopics)},
	}
}

func allTopicsWithCoverageCallbackUpdate(data string) *tgbotapi.Update {
	return &tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Data: transport.AllTopicsWithCoverage,
		},
	}
}

func createCampaignCallbackUpdate(data string) *tgbotapi.Update {
	return &tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Data: transport.CreateCampaign,
		},
	}
}

func createCampaignMessageUpdate(data string) *tgbotapi.Update {
	return &tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Text: "Food",
		},
	}
}

func myCampaignsCallbackUpdate(data string) *tgbotapi.Update {
	return &tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Data: transport.MyCampaigns,
		},
	}
}

func getCampaignID(db *gorm.DB) string {
	var campaign models.Campaign
	db.First(&campaign)

	return campaign.ID.String()
}

func createAdCallbackUpdate(campaignID string) *tgbotapi.Update {
	fmt.Println("### campaignID", campaignID)
	return &tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Data: fmt.Sprintf("%s/%s", transport.CreateAd, campaignID),
		},
	}
}

func createAdMessageUpdate(data string) *tgbotapi.Update {
	return &tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Text: `
Name: Stock market
TargetTopics: art, food
BudgetUSD: 100
CostPerView: 0.1
Message: Follow this [link](https://www.investing.com/) to find more about investments!`,
		},
	}
}
