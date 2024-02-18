package service

import (
	"advertiser/shared/pkg/repo/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UseCase interface {
	AllTopics() (res []string)
	StoreInitialChannelData([]tgbotapi.ChatMember, models.Channel) error
	DeleteChannel(chatID int64) error

	ListMyChannels(userID int64) (map[int64]string, error)
	GetChannelInfo(channelID int64) (channel *models.Channel, err error)
	UpdateChannelTopics(channelID int64, topics []string) (err error)
}
