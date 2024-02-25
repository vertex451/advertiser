package listener

import (
	"advertiser/channel_owner/internal/service/listener/repo/postgresql/types"
	"advertiser/shared/pkg/service/repo/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UseCase interface {
	Listener
	Writer
}

type Listener interface {
	AllTopics() (res []string)
	StoreInitialChannelData([]tgbotapi.ChatMember, models.Channel) error
	DeleteChannel(chatID int64) error

	ListMyChannels(userID int64) (map[int64]string, error)
	GetChannelInfo(channelID int64) (channel *models.Channel, err error)
	UpdateChannelTopics(channelID int64, topics []string) (err error)
}

type Writer interface {
	CheckForNewAds() ([]types.GetAdsOnModerationResult, error)
}
