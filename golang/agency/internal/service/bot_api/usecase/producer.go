package usecase

import (
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/pkg/service/usecase"
	uuid "github.com/satori/go.uuid"
	"slices"
	"sort"
	"tg-bot/internal/service/bot_api/types"
)

func (uc *UseCase) AllTopics() (res []string) {
	for topic := range uc.cache.topics {
		res = append(res, topic)
	}

	slices.Sort(res)

	return res
}

func (uc *UseCase) AllTopicsWithCoverage() (res []types.TopicWithCoverage, err error) {
	topicsMap := uc.cache.topics
	for topic, coverage := range topicsMap {
		res = append(res, types.TopicWithCoverage{
			Name:     topic,
			Coverage: coverage,
		})
	}

	// first, sort by name
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	// second, sort by coverage
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Coverage > res[j].Coverage
	})

	return res, nil
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

func (uc *UseCase) UpsertAd(advertisement *models.Advertisement) error {
	var topics []string
	for _, topic := range advertisement.TargetTopics {
		topics = append(topics, topic.ID)
	}
	err := usecase.ValidateTopics(uc.topics, topics)
	if err != nil {
		return err
	}

	return uc.repo.UpsertAd(advertisement)
}

func (uc *UseCase) GetAdDetails(id uuid.UUID) (*models.Advertisement, error) {
	return uc.repo.GetAdDetails(id)
}

func (uc *UseCase) RunAd(id uuid.UUID) error {
	return uc.repo.RunAd(id)
}
