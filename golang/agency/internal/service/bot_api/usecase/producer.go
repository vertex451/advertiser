package usecase

import (
	"advertiser/shared/pkg/repo/models"
	"advertiser/shared/pkg/utils"

	uuid "github.com/satori/go.uuid"
	"tg-bot/internal/service/bot_api/usecase/types"
)

func (uc *UseCase) AllTopics() (res []string) {
	for topic := range uc.cache.topics {
		res = append(res, topic)
	}

	return res
}

func (uc *UseCase) AllTopicsWithCoverage() (map[string]int, error) {
	return uc.cache.topics, nil
}

func (uc *UseCase) CreateCampaign(userID int64, campaignName string) (uuid.UUID, error) {
	return uc.repo.CreateCampaign(userID, campaignName)
}

func (uc *UseCase) ListMyCampaigns(userID int64) ([]models.Campaign, error) {
	return uc.repo.ListMyCampaigns(userID)
}

func (uc *UseCase) CampaignDetails(campaignID uuid.UUID) (*models.Campaign, error) {
	return uc.repo.CampaignDetails(campaignID)
}

func (uc *UseCase) UpsertAd(advertisement models.Advertisement) (*uuid.UUID, error) {
	var topics []string
	for _, topic := range advertisement.TargetTopics {
		topics = append(topics, topic.ID)
	}
	err := utils.ValidateTopics(uc.topics, topics)
	if err != nil {
		return nil, err
	}

	return uc.repo.UpsertAd(advertisement)
}

func (uc *UseCase) GetAdDetails(id uuid.UUID) (*types.Advertisement, error) {
	ad, err := uc.repo.GetAdDetails(id)
	if err != nil {
		return nil, err
	}

	coverage := 0
	var topicList []string
	for _, topic := range ad.TargetTopics {
		topicList = append(topicList, topic.ID)
		coverage += uc.cache.topics[topic.ID]
	}

	return &types.Advertisement{
		ID:       ad.ID,
		Name:     ad.Name,
		Topics:   topicList,
		Coverage: coverage,
		Budget:   ad.Budget,
		Message:  ad.Message,
	}, nil
}

func (uc *UseCase) EditAd(ad models.Advertisement) (*models.Advertisement, error) {
	return uc.repo.EditAd(ad)
}
