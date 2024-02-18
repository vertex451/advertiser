package usecase

import (
	"advertiser/shared/pkg/repo/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (uc *UseCase) AllTopics() (res []string) {
	for topic := range uc.cache.topics {
		res = append(res, topic)
	}

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
	err = uc.validateTopics(topics)
	if err != nil {
		return err
	}

	err = uc.repo.UpdateChannelTopics(channelID, topics)
	if err != nil {
		return err
	}

	return uc.updateTopicCache()
}

func (uc *UseCase) DeleteChannel(chatID int64) error {
	return uc.repo.DeleteChannel(chatID)
}
