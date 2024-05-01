package tg_bot_api

import (
	"advertiser/shared/utils"
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

func (t *TgBotApi) GetToken() string {
	return t.tgBotApi.Token
}

func (t *TgBotApi) GetUpdatesChan() tgbotapi.UpdatesChannel {
	return t.tgBotApi.GetUpdatesChan(t.updateConfig)
}

func (t *TgBotApi) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	// Use RetryWithBackoff from the utils package to handle retries and backoff
	result, err := utils.RetryWithBackoff(func() (interface{}, error) {
		return t.tgBotApi.Send(c)
	}, utils.DefaultMaxRetries, utils.DefaultBaseBackoffTime)

	if err != nil {
		// Return an empty message and the error if all attempts failed
		return tgbotapi.Message{}, err
	}

	// Type assert the result to the expected return type
	return result.(tgbotapi.Message), nil
}

// GetChatAdministrators with retry and backoff
func (t *TgBotApi) GetChatAdministrators(config tgbotapi.ChatAdministratorsConfig) ([]tgbotapi.ChatMember, error) {
	result, err := utils.RetryWithBackoff(func() (interface{}, error) {
		return t.tgBotApi.GetChatAdministrators(config)
	}, utils.DefaultMaxRetries, utils.DefaultBaseBackoffTime)

	if err != nil {
		return nil, err
	}
	return result.([]tgbotapi.ChatMember), nil
}

// GetChatMembersCount with retry and backoff
func (t *TgBotApi) GetChatMembersCount(config tgbotapi.ChatMemberCountConfig) (int, error) {
	result, err := utils.RetryWithBackoff(func() (interface{}, error) {
		return t.tgBotApi.GetChatMembersCount(config)
	}, utils.DefaultMaxRetries, utils.DefaultBaseBackoffTime)

	// Type assert the result to the expected return type
	if err != nil {
		return 0, err
	}
	return result.(int), nil
}

func (t *TgBotApi) GetFile(cfg tgbotapi.FileConfig) (tgbotapi.File, error) {
	// Use RetryWithBackoff from the utils package to handle retries and backoff
	result, err := utils.RetryWithBackoff(func() (interface{}, error) {
		return t.tgBotApi.GetFile(cfg)
	}, utils.DefaultMaxRetries, utils.DefaultBaseBackoffTime)

	if err != nil {
		// Return an empty file and the error if all attempts failed
		return tgbotapi.File{}, err
	}

	// Type assert the result to the expected return type
	return result.(tgbotapi.File), nil
}

func (t *TgBotApi) GetFileDirectURL(filePath string) (string, error) {
	result, err := utils.RetryWithBackoff(func() (interface{}, error) {
		return t.tgBotApi.GetFileDirectURL(filePath)
	}, utils.DefaultMaxRetries, utils.DefaultBaseBackoffTime)

	if err != nil {
		return "", err
	}

	return result.(string), nil
}
