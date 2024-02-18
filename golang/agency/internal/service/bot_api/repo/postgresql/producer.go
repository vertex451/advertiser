package postgresql

import (
	"advertiser/shared/pkg/repo/models"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"sort"
)

func (r *Repository) AllTopics() (res []string, err error) {
	var allTopics []*models.Topic
	if err = r.Db.Find(&allTopics).Error; err != nil {
		return nil, err
	}

	for _, topic := range allTopics {
		res = append(res, topic.ID)
	}

	return res, nil
}

func (r *Repository) AllTopicsWithCoverage() (map[string]int, error) {
	res := make(map[string]int)
	var topicsWithSubscribers []struct {
		Name        string
		Subscribers int
	}

	err := r.Db.Raw(`
SELECT t.id as name, SUM(c.subscribers) as subscribers
FROM topics AS t
         LEFT JOIN
     channel_topics as CT on t.id = ct.topic_id
         LEFT JOIN channels c on CT.channel_id = c.id
GROUP BY t.id;
`).Find(&topicsWithSubscribers).Error
	if err != nil {
		zap.L().Error("failed to get topics", zap.Error(err))
		return nil, err
	}

	for _, ts := range topicsWithSubscribers {
		res[ts.Name] = ts.Subscribers
	}

	return res, nil
}

func (r *Repository) CreateCampaign(userID int64, campaignName string) (uuid.UUID, error) {
	campaign := models.Campaign{
		Name:   campaignName,
		UserID: userID,
	}

	err := r.Db.Create(&campaign).Error

	return campaign.ID, err
}

func (r *Repository) ListMyCampaigns(userID int64) ([]models.Campaign, error) {
	var campaigns []models.Campaign
	err := r.Db.Where("user_id = ?", userID).Find(&campaigns).Error
	if err != nil {
		return nil, err
	}

	sort.SliceStable(campaigns, func(i, j int) bool {
		return campaigns[i].Name < campaigns[j].Name
	})

	return campaigns, nil
}

func (r *Repository) CampaignDetails(campaignID uuid.UUID) (*models.Campaign, error) {
	campaign := models.Campaign{
		ID: campaignID,
	}

	err := r.Db.Preload("Advertisements").Find(&campaign).Error
	if err != nil {
		return nil, err
	}

	return &campaign, err
}

func (r *Repository) UpsertAd(advertisement models.Advertisement) (*uuid.UUID, error) {
	var err error
	if advertisement.ID != uuid.Nil {
		err = r.Db.Save(&advertisement).Error
	} else {
		err = r.Db.Create(&advertisement).Error
	}

	return &advertisement.ID, err
}

func (r *Repository) GetAdDetails(id uuid.UUID) (*models.Advertisement, error) {
	ad := models.Advertisement{
		ID: id,
	}

	err := r.Db.Find(&ad).Error
	if err != nil {
		return nil, err
	}

	return &ad, err
}

func (r *Repository) EditAd(ad models.Advertisement) (*models.Advertisement, error) {
	if ad.ID == uuid.Nil {
		zap.L().Panic("advertisement uuid is nil!")
	}

	err := r.Db.Save(&ad).Error
	if err != nil {
		return nil, err
	}

	return &ad, err
}
