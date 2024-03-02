package postgresql

import (
	"advertiser/shared/pkg/service/repo/models"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

func (r *Repository) GetAdsOnModeration() (res []models.AdvertisementChannel, err error) {
	err = r.Db.Raw(`
SELECT c.id as channel_id, c.title as channel_title, 
       ca.user_id as channel_owner_id, 
       ads.id as advertisement_id, ads.name as ad_name,  
       ads.message as ad_message, ads.cost_per_view as ad_cost_per_view
FROM advertisements as ads
         LEFT JOIN advertisement_topics at on ads.id = at.advertisement_id
         LEFT JOIN channel_topics ct on at.topic_id = ct.topic_id
         LEFT JOIN channels c on ct.channel_id = c.id
         LEFT JOIN channel_admins ca on c.id = ca.channel_id
WHERE ads.status = 'pending' AND ca.role = 'creator'
GROUP BY c.id, ca.user_id, ads.message, ads.id;
`).Find(&res).Error
	if err != nil {
		zap.L().Error("failed to get topics", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (r *Repository) CreateAdvertisementChannelEntries(ads []models.AdvertisementChannel) {
	var err error
	success := make(map[uuid.UUID]struct{})
	for _, entry := range ads {
		entry.Status = models.AdChanCreated
		err = r.Db.Create(&entry).Error
		if err != nil {
			zap.L().Error("failed to create advertisement channel entry", zap.Error(err))
			continue
		}
		success[entry.AdvertisementID] = struct{}{}
	}

	for id := range success {
		err = r.Db.Model(&models.Advertisement{}).
			Where("id = ?", id).
			Update("status", models.AdsStatusRunning).Error
		if err != nil {
			zap.L().Error("failed to update advertisement status", zap.Error(err))
		}
	}
}

func (r *Repository) GetAdsChannelByStatus(status models.AdChanStatus) (res []models.AdvertisementChannel, err error) {
	err = r.Db.Where("status = ?", status).Find(&res).Error
	if err != nil {
		zap.L().Error("failed to get advertisement channel by status", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (r *Repository) UpdateAdChanStatus(adChannelID string, status models.AdChanStatus) error {
	return r.Db.Model(&models.AdvertisementChannel{}).
		Where("id = ?", adChannelID).
		Update("status", status).Error
}

func (r *Repository) SetAdChanMessageID(adChanID string, msgID int) error {
	return r.Db.Model(&models.AdvertisementChannel{}).
		Where("id = ?", adChanID).
		Update("message_id", msgID).Error
}
