package postgresql

import (
	"advertiser/channel_owner/internal/service/listener/repo/postgresql/types"
	"go.uber.org/zap"
)

func (r *Repository) GetAdsOnModeration() (res []types.GetAdsOnModerationResult, err error) {
	err = r.Db.Raw(`
SELECT c.id as channel_id, c.title, ca.user_id, ads.message as ad_message, ads.id as ad_id
FROM advertisements as ads
         LEFT JOIN advertisement_topics at on ads.id = at.advertisement_id
         LEFT JOIN channel_topics ct on at.topic_id = ct.topic_id
         LEFT JOIN channels c on ct.channel_id = c.id
         LEFT JOIN channel_admins ca on c.id = ca.channel_id
WHERE ads.status = 'pending_approval' AND ca.role = 'creator'
GROUP BY c.id, ca.user_id, ads.message, ads.id;
`).Find(&res).Error
	if err != nil {
		zap.L().Error("failed to get topics", zap.Error(err))
		return nil, err
	}

	return res, nil
}
