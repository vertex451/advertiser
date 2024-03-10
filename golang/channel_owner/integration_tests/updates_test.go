package integration_tests

import (
	"advertiser/channel_owner/internal/service/listener/transport"
	"advertiser/channel_owner/internal/service/listener/transport/mocks"
	"advertiser/shared/pkg/service/constants"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type channel struct {
	id int64
}

type bot struct {
	id int64
}

type testCase struct {
	testName        string
	update          tgbotapi.Update
	expectedMsgText string
	expectedButtons []tgbotapi.InlineKeyboardButton
}

var (
	testChannel = channel{id: -1002093237940}
	testBot     = bot{id: 6406834985}
)

func startCommandUpdate() tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			From:     &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: len(constants.Start) + 1}},
			Text:     fmt.Sprintf("/%s", constants.Start)},
	}
}

func startCallbackUpdate() tgbotapi.Update {
	return tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Data: constants.Start,
		},
	}
}

func allTopicsCommandUpdate() tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			From:     &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: len(constants.AllTopics) + 1}},
			Text:     fmt.Sprintf("/%s", constants.AllTopics)},
	}
}

func allTopicsCallbackUpdate() tgbotapi.Update {
	return tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Data: constants.AllTopics,
		},
	}
}

func botIsAddedToChannelUpdate() tgbotapi.Update {
	return tgbotapi.Update{
		MyChatMember: &tgbotapi.ChatMemberUpdated{
			Chat: tgbotapi.Chat{
				ID:       testChannel.id,
				Type:     "channel",
				Title:    "Sport Channel",
				UserName: "sport_channel123",
			},
			NewChatMember: tgbotapi.ChatMember{
				User: &tgbotapi.User{
					ID:       testBot.id,
					UserName: transport.ChannelMonetizerBotName,
				},
				Status:            transport.StatusAdministrator,
				CanPostMessages:   true,
				CanDeleteMessages: true,
			},
		},
	}
}

func myChannelsCallbackUpdate() tgbotapi.Update {
	return tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Data: transport.MyChannels,
		},
	}
}

func editTopicsCallbackUpdate() tgbotapi.Update {
	return tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Data: fmt.Sprintf("%s/%d", transport.EditChannelsTopics, testChannel.id),
		},
	}
}

func editTopicsMessageUpdate() tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			From: &tgbotapi.User{ID: mocks.ChannelCreator.ID},
			Text: "art, books, food, pets, sport",
		},
	}
}

func botIsRemovedFromChannelUpdate() tgbotapi.Update {
	return tgbotapi.Update{
		MyChatMember: &tgbotapi.ChatMemberUpdated{
			Chat: tgbotapi.Chat{
				ID:       testChannel.id,
				Type:     "channel",
				Title:    "Sport Channel",
				UserName: "sport_channel123",
			},
			NewChatMember: tgbotapi.ChatMember{
				User: &tgbotapi.User{
					ID:       testBot.id,
					UserName: transport.ChannelMonetizerBotName,
				},
				Status: transport.StatusLeft,
			},
		},
	}
}
