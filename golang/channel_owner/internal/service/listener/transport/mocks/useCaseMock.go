package mocks

import (
	"advertiser/shared/pkg/service/repo/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UseCaseMock struct {
}

func (u UseCaseMock) AllTopics() (res []string) {
	return []string{"topic1", "topic2", "topic3"}
}

func (u UseCaseMock) StoreInitialChannelData(members []tgbotapi.ChatMember, channel models.Channel) error {
	//TODO implement me
	panic("implement me")
}

func (u UseCaseMock) DeleteChannel(chatID int64) error {
	//TODO implement me
	panic("implement me")
}

func (u UseCaseMock) ListMyChannels(userID int64) (map[int64]string, error) {
	//TODO implement me
	panic("implement me")
}

func (u UseCaseMock) GetChannelInfo(channelID int64) (channel *models.Channel, err error) {
	//TODO implement me
	panic("implement me")
}

func (u UseCaseMock) UpdateChannelTopics(channelID int64, topics []string) (err error) {
	//TODO implement me
	panic("implement me")
}

func (u UseCaseMock) GetAdsToModerateByUserID(id int64) ([]models.AdvertisementChannel, error) {
	//TODO implement me
	panic("implement me")
}

func (u UseCaseMock) GetAdChanDetails(id string) (*models.AdvertisementChannel, error) {
	//TODO implement me
	panic("implement me")
}

func (u UseCaseMock) UpdateAdChanEntry(channel models.AdvertisementChannel) error {
	//TODO implement me
	panic("implement me")
}

func (u UseCaseMock) UpdateAd(ad models.Advertisement) error {
	//TODO implement me
	panic("implement me")
}

func (u UseCaseMock) CheckForNewAds() {
	//TODO implement me
	panic("implement me")
}

func (u UseCaseMock) GetAdsChannelByStatus(status models.AdChanStatus) (map[int64][]models.AdvertisementChannel, error) {
	//TODO implement me
	panic("implement me")
}

func (u UseCaseMock) GetRunningAdvertisements() ([]*models.Advertisement, error) {
	//TODO implement me
	panic("implement me")
}
