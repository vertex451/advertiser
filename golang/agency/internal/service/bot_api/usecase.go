package bot_api

import (
	"advertiser/shared/pkg/service/repo/models"
	uuid "github.com/satori/go.uuid"
	"tg-bot/internal/service/bot_api/types"
)

type UseCase interface {
	AllTopics() (res []string)
	AllTopicsWithCoverage() (map[string]int, error)
	CreateCampaign(respondTo int64, campaignName string) (uuid.UUID, error)
	ListMyCampaigns(userID int64) ([]models.Campaign, error)
	CampaignDetails(campaignID uuid.UUID) (*models.Campaign, error)

	UpsertAd(advertisement models.Advertisement) (*uuid.UUID, error)
	GetAdDetails(uuid.UUID) (*types.Advertisement, error)
	RunAd(uuid.UUID) error
}
