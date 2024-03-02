package listener

import (
	"advertiser/shared/pkg/service/repo/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Repo interface {
	RepoListener
	RepoNotification
}

type RepoListener interface {
	AllTopics() ([]string, error)
	StoreInitialChannelData(admins []tgbotapi.ChatMember, chat models.Channel) error
	DeleteChannel(chatID int64) error

	ListMyChannels(userID int64) (map[int64]string, error)
	GetChannelInfo(channelID int64) (channel *models.Channel, err error)
	UpdateChannelTopics(channelID int64, topics []string) (err error)

	GetAdsToModerateByUserID(id int64) ([]models.AdvertisementChannel, error)

	GetAdChanDetails(id string) (*models.AdvertisementChannel, error)
	SetAdChanMessageID(adChanID string, msgID int) error
}

type RepoNotification interface {
	GetAdsOnModeration() (res []models.AdvertisementChannel, err error)
	CreateAdvertisementChannelEntries(ads []models.AdvertisementChannel)
	GetAdsChannelByStatus(status models.AdChanStatus) (res []models.AdvertisementChannel, err error)
	UpdateAdChanStatus(adChannelID string, status models.AdChanStatus) error
}
