package bot_api

import (
	"advertiser/shared/pkg/service/repo/models"
	uuid "github.com/satori/go.uuid"
	"tg-bot/internal/service/bot_api/types"
)

type UseCase interface {
	AllTopics() (res []string)
	AllTopicsWithCoverage() ([]types.TopicWithCoverage, error)
	CreateCampaign(respondTo int64, campaignName string) (uuid.UUID, error)
	ListMyCampaigns(userID int64) ([]models.Campaign, error)
	CampaignDetails(campaignID uuid.UUID) (*models.Campaign, error)

	UpsertAd(advertisement *models.Advertisement) error
	GetAdDetails(uuid.UUID) (*models.Advertisement, error)
	RunAd(uuid.UUID) error
}
