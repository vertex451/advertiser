package listener

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/repo/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	uuid "github.com/satori/go.uuid"
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
	UpdateChannelCostPerMile(channelID int64, costPerMile float64) (err error)
	UpdateChannelLocation(channelID int64, location constants.Location) (err error)

	GetAdsToModerateByUserID(id int64) ([]models.AdvertisementChannel, error)

	GetAdChanDetails(id string) (*models.AdvertisementChannel, error)
	GetAdMessageByAdChanID(uuid.UUID) (*models.Advertisement, error)

	UpdateAdChanEntry(channel models.AdvertisementChannel) error
	UpdateAd(ad models.Advertisement) error

	ReportBug(userID int64, message string) error
	RequestFeature(userID int64, message string) error
}

type Writer interface {
	CheckForNewAds()
	GetAdsChannelByStatus(status models.AdChanStatus) (map[int64][]models.AdvertisementChannel, error)
	GetRunningAdvertisements() ([]*models.Advertisement, error)
}
