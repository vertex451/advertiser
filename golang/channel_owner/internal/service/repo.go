package service

import (
	"advertiser/shared/pkg/repo/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Repo interface {
	AllTopics() ([]string, error)
	StoreInitialChannelData(admins []tgbotapi.ChatMember, chat models.Channel) error
	DeleteChannel(chatID int64) error

	ListMyChannels(userID int64) (map[int64]string, error)
	GetChannelInfo(channelID int64) (channel *models.Channel, err error)
	UpdateChannelTopics(channelID int64, topics []string) (err error)
}
