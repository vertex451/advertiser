package mocks

import (
	"advertiser/channel_owner/internal/service/listener/transport"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBotApiMock struct {
	UpdatesChan tgbotapi.UpdatesChannel
	TargetChan  chan tgbotapi.Chattable
}

type User struct {
	ID       int64
	UserName string
}

var ChannelCreator = User{
	ID:       1,
	UserName: "thisIsCreator",
}

func (t TgBotApiMock) GetChatAdministrators(config tgbotapi.ChatAdministratorsConfig) ([]tgbotapi.ChatMember, error) {
	return []tgbotapi.ChatMember{
		{
			User: &tgbotapi.User{
				ID:       ChannelCreator.ID,
				UserName: ChannelCreator.UserName,
			},
			Status: transport.StatusCreator,
		},
	}, nil
}

func (t TgBotApiMock) GetChatMembersCount(config tgbotapi.ChatMemberCountConfig) (int, error) {
	return 1000, nil
}

func (t TgBotApiMock) GetUpdatesChan() tgbotapi.UpdatesChannel {
	return t.UpdatesChan
}

func (t TgBotApiMock) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	//TODO implement me
	t.TargetChan <- c

	return tgbotapi.Message{}, nil
}
