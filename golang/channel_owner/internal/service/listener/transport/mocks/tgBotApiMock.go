package mocks

import (
	"advertiser/channel_owner/internal/service/listener/transport"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBotApiMock struct {
	UpdatesChan tgbotapi.UpdatesChannel
	TargetChan  chan tgbotapi.Chattable
}

func (t TgBotApiMock) GetChatAdministrators(config tgbotapi.ChatAdministratorsConfig) ([]tgbotapi.ChatMember, error) {
	return []tgbotapi.ChatMember{
		{
			User: &tgbotapi.User{
				ID:       1,
				UserName: "thisIsCreator",
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
