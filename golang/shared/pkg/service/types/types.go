package types

import (
	"advertiser/shared/pkg/service/repo/models"
	uuid "github.com/satori/go.uuid"
)

type GetNewAdsCountByUserIDResponse struct {
	AdvertisementId uuid.UUID
	ChannelId       int64
	UserId          int64
	Count           int64
}

// Advertisement is needed because of additional calculation we do in UseCase.getAdDetails
type Advertisement struct {
	ID uuid.UUID

	// metadata
	Name        string
	Topics      []string
	Coverage    int
	Budget      int
	MsgText     string
	MsgEntities []models.MsgEntity
	MsgImageURL string
	Status      models.AdvertisementStatus
}
