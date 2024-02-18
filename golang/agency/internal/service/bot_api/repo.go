package bot_api

import (
	"advertiser/shared/pkg/repo/models"
	uuid "github.com/satori/go.uuid"
)

type Repo interface {
	AllTopics() ([]string, error)
	AllTopicsWithCoverage() (map[string]int, error)

	CreateCampaign(userID int64, campaignName string) (uuid.UUID, error)
	ListMyCampaigns(userID int64) ([]models.Campaign, error)
	CampaignDetails(campaignID uuid.UUID) (*models.Campaign, error)

	UpsertAd(advertisement models.Advertisement) (*uuid.UUID, error)
	GetAdDetails(id uuid.UUID) (*models.Advertisement, error)
	EditAd(models.Advertisement) (*models.Advertisement, error)
}
