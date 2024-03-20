package tg_bot_api

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TgBotApiProvider interface {
	GetChatAdministrators(config tgbotapi.ChatAdministratorsConfig) ([]tgbotapi.ChatMember, error)
	GetChatMembersCount(config tgbotapi.ChatMemberCountConfig) (int, error)
	GetFile(cfg tgbotapi.FileConfig) (tgbotapi.File, error)
	GetToken() string
	GetFileDirectURL(filePath string) (string, error)
	GetUpdatesChan() tgbotapi.UpdatesChannel
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type TgBotApi struct {
	tgBotApi     *tgbotapi.BotAPI
	updateConfig tgbotapi.UpdateConfig
}

func New(token string) *TgBotApi {
	if token == "" {
		zap.L().Panic("token is empty")
	}
	tgBotApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		zap.L().Panic("failed to initiate tgBotApi", zap.Error(err))
	}
	tgBotApi.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	return &TgBotApi{
		tgBotApi:     tgBotApi,
		updateConfig: updateConfig,
	}
}

func (t *TgBotApi) GetUpdatesChan() tgbotapi.UpdatesChannel {
	return t.tgBotApi.GetUpdatesChan(t.updateConfig)
}

func (t *TgBotApi) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return t.tgBotApi.Send(c)
}

func (t *TgBotApi) GetChatAdministrators(config tgbotapi.ChatAdministratorsConfig) ([]tgbotapi.ChatMember, error) {
	return t.tgBotApi.GetChatAdministrators(config)
}

func (t *TgBotApi) GetChatMembersCount(config tgbotapi.ChatMemberCountConfig) (int, error) {
	return t.tgBotApi.GetChatMembersCount(config)
}

func (t *TgBotApi) GetFile(cfg tgbotapi.FileConfig) (tgbotapi.File, error) {
	return t.tgBotApi.GetFile(cfg)
}

func (t *TgBotApi) GetFileDirectURL(filePath string) (string, error) {
	return t.tgBotApi.GetFileDirectURL(filePath)
}

func (t *TgBotApi) GetToken() string {
	return t.tgBotApi.Token
}
