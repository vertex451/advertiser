package integration_tests

import (
	"advertiser/channel_owner/internal/service/listener/transport"
	"advertiser/shared/pkg/service/constants"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

type channel struct {
	id int64
}

type bot struct {
	chatID int64
}

type testCase struct {
	testName        string
	update          tgbotapi.Update
	expectedMsgText string
	expectedButtons []tgbotapi.InlineKeyboardButton
}

var (
	testChannel = channel{id: -1002093237940}
	testBot     = bot{chatID: 6406834985}
)

func TestStart(t *testing.T) {
	updatesChan := make(chan tgbotapi.Update)
	targetChan := make(chan tgbotapi.Chattable)

	go startChannelOwnerService(updatesChan, targetChan)

	tc := []testCase{
		{
			testName: "TestStartCommand",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: len(constants.Start) + 1}},
					Chat:     &tgbotapi.Chat{ID: 6406834985},
					Text:     fmt.Sprintf("/%s", constants.Start)},
			},
			expectedMsgText: "Choose action:",
		},
		{
			testName: "TestStartCallback",
			update: tgbotapi.Update{
				CallbackQuery: &tgbotapi.CallbackQuery{
					Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: testBot.chatID}}, Data: constants.Start,
				},
			},
			expectedMsgText: "Choose action:",
		},
	}

	for _, tt := range tc {
		t.Run(tt.testName, func(t *testing.T) {
			updatesChan <- tt.update
			msgRaw, _ := <-targetChan
			msg := msgRaw.(tgbotapi.MessageConfig)
			assert.Equal(t, tt.expectedMsgText, msg.Text)
		})
	}
}

func TestAllTopicsCommand(t *testing.T) {
	updatesChan := make(chan tgbotapi.Update)
	targetChan := make(chan tgbotapi.Chattable)

	go startChannelOwnerService(updatesChan, targetChan)

	tc := []testCase{
		{
			testName: "TestAllTopicsCommand",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: len(constants.AllTopics) + 1}},
					Chat:     &tgbotapi.Chat{ID: testBot.chatID},
					Text:     fmt.Sprintf("/%s", constants.AllTopics)},
			},
			expectedMsgText: `
Supported topics:
art, books, food, pets, sport
`,
		},
		{
			testName: "TestAllTopicsCallback",
			update: tgbotapi.Update{
				CallbackQuery: &tgbotapi.CallbackQuery{
					Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 6406834985}}, Data: constants.AllTopics,
				},
			},
			expectedMsgText: `
Supported topics:
art, books, food, pets, sport
`,
		},
	}

	for _, tt := range tc {
		t.Run(tt.testName, func(t *testing.T) {
			updatesChan <- tt.update
			msgRaw, _ := <-targetChan
			msg := msgRaw.(tgbotapi.MessageConfig)
			assert.Equal(t, tt.expectedMsgText, msg.Text)
		})
	}
}

// TestChannelOwnerFlow
// 1. Handle event of adding bot to channel
// 2. Set topics
// 3. List my channels
func TestChannelOwnerFlow(t *testing.T) {
	updatesChan := make(chan tgbotapi.Update)
	targetChan := make(chan tgbotapi.Chattable)

	go startChannelOwnerService(updatesChan, targetChan)

	tc := []testCase{
		{
			testName: "TestBotIsAddedToChannelEvent",
			update: tgbotapi.Update{
				MyChatMember: &tgbotapi.ChatMemberUpdated{
					Chat: tgbotapi.Chat{
						ID:       testChannel.id,
						Type:     "channel",
						Title:    "Sport Channel",
						UserName: "sport_channel123",
					},
					NewChatMember: tgbotapi.ChatMember{
						User: &tgbotapi.User{
							ID:       testBot.chatID,
							UserName: transport.ChannelMonetizerBotName,
						},
						Status:            transport.StatusAdministrator,
						CanPostMessages:   true,
						CanDeleteMessages: true,
					},
				},
			},
			expectedMsgText: "Advertiser bot was successfully added to Sport Channel",
		},
		{
			testName: "TestMyChannelsCallback",
			update: tgbotapi.Update{
				CallbackQuery: &tgbotapi.CallbackQuery{
					Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: testBot.chatID}},
					Data:    transport.MyChannels,
				},
			},
			expectedMsgText: "Select a channel:",
		},
		{
			testName: "TestEditTopicsCallback",
			update: tgbotapi.Update{
				CallbackQuery: &tgbotapi.CallbackQuery{
					Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: testBot.chatID}},
					Data:    fmt.Sprintf("%s/%d", transport.EditChannelsTopics, testChannel.id),
				},
			},
			expectedMsgText: `
Choose topics from the list:
art, books, food, pets, sport
`,
		},
		{
			// TestEditTopicsMessage must go after TestEditTopicsCallback due to state nature
			testName: "TestEditTopicsMessage",
			update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					Chat: &tgbotapi.Chat{ID: testBot.chatID},
					Text: "art, books, food, pets, sport",
				},
			},
			expectedMsgText: "Topics changed! New channel topics: art, books, food, pets, sport",
		},
		{
			testName: "TestBotIsRemovedFromChannelEvent",
			update: tgbotapi.Update{
				MyChatMember: &tgbotapi.ChatMemberUpdated{
					Chat: tgbotapi.Chat{
						ID:       testChannel.id,
						Type:     "channel",
						Title:    "Sport Channel",
						UserName: "sport_channel123",
					},
					OldChatMember: tgbotapi.ChatMember{
						User: &tgbotapi.User{
							ID:       testBot.chatID,
							UserName: transport.ChannelMonetizerBotName,
						},
						Status: transport.StatusLeft,
					},
				},
			},
			expectedMsgText: "Advertiser bot was successfully added to Sport Channel",
		},
	}

	for _, tt := range tc {
		t.Run(tt.testName, func(t *testing.T) {
			updatesChan <- tt.update
			msgRaw, _ := <-targetChan
			msg := msgRaw.(tgbotapi.MessageConfig)
			assert.Equal(t, tt.expectedMsgText, msg.Text)
		})
	}
}
