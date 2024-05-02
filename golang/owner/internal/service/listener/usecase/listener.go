package usecase

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/pkg/service/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	uuid "github.com/satori/go.uuid"
	"slices"
)

func (uc *UseCase) AllTopics() (res []string) {
	for topic := range uc.cache.topics {
		res = append(res, topic)
	}

	slices.Sort(res)

	return res
}

func (uc *UseCase) StoreInitialChannelData(admins []tgbotapi.ChatMember, chat models.Channel) error {
	return uc.repo.StoreInitialChannelData(admins, chat)
}

func (uc *UseCase) ListMyChannels(userID int64) (map[int64]string, error) {
	return uc.repo.ListMyChannels(userID)
}

func (uc *UseCase) GetChannelInfo(channelID int64) (channel *models.Channel, err error) {
	return uc.repo.GetChannelInfo(channelID)
}

func (uc *UseCase) UpdateChannelTopics(channelID int64, topics []string) (err error) {
	err = usecase.ValidateTopics(uc.topics, topics)
	if err != nil {
		return err
	}

	err = uc.repo.UpdateChannelTopics(channelID, topics)
	if err != nil {
		return err
	}

	return uc.updateTopicCache()
}

func (uc *UseCase) UpdateChannelLocation(channelID int64, location constants.Location) (err error) {
	return uc.repo.UpdateChannelLocation(channelID, location)
}

func (uc *UseCase) DeleteChannel(chatID int64) error {
	return uc.repo.DeleteChannel(chatID)
}

func (uc *UseCase) GetAdsToModerateByUserID(id int64) ([]models.AdvertisementChannel, error) {
	return uc.repo.GetAdsToModerateByUserID(id)
}

func (uc *UseCase) GetAdChanDetails(id string) (*models.AdvertisementChannel, error) {
	return uc.repo.GetAdChanDetails(id)
}

func (uc *UseCase) GetAdMessageByAdChanID(id uuid.UUID) (*models.Advertisement, error) {
	return uc.repo.GetAdMessageByAdChanID(id)
}

func (uc *UseCase) ReportBug(userID int64, message string) error {
	return uc.repo.ReportBug(userID, message)
}

func (uc *UseCase) RequestFeature(userID int64, message string) error {
	return uc.repo.RequestFeature(userID, message)
}
