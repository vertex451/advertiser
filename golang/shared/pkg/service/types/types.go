package types

import uuid "github.com/satori/go.uuid"

type GetNewAdsCountByUserIDResponse struct {
	AdvertisementId uuid.UUID
	ChannelId       int64
	UserId          int64
	Count           int64
}
