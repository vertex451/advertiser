package listener

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/repo/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	uuid "github.com/satori/go.uuid"
)

type Repo interface {
	AllTopics() ([]string, error)

	Channel
	Advertisement
}

type Channel interface {
	StoreInitialChannelData(admins []tgbotapi.ChatMember, chat models.Channel) error
	UpdateChannelTopics(channelID int64, topics []string) (err error)
	UpdateChannelCostPerMile(channelID int64, costPerMile float64) (err error)
	UpdateChannelLocation(channelID int64, location constants.Location) (err error)

	GetChannelInfo(channelID int64) (channel *models.Channel, err error)
	ListMyChannels(userID int64) (map[int64]string, error)
	DeleteChannel(chatID int64) error
}

type Advertisement interface {
	UpdateAd(ad models.Advertisement) error
	GetAdsOnModeration() (res []models.AdvertisementChannel, err error)
	GetAdsToModerateByUserID(id int64) ([]models.AdvertisementChannel, error)
	GetRunningAds() ([]*models.Advertisement, error)

	CreateAdChanEntries(ads []models.AdvertisementChannel)
	UpdateAdChanEntry(channel models.AdvertisementChannel) error
	GetAdChanDetails(id string) (*models.AdvertisementChannel, error)
	GetAdMessageByAdChanID(id uuid.UUID) (*models.Advertisement, error)
	GetAdChannelByStatus(status models.AdChanStatus) (res []models.AdvertisementChannel, err error)

	ReportBug(userID int64, message string) error
	RequestFeature(userID int64, message string) error
}
