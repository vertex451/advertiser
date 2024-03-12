package types

import (
	"advertiser/shared/pkg/service/repo/models"
	uuid "github.com/satori/go.uuid"
)

// Advertisement is needed because of additional calculation we do in UseCase.GetAdDetails
type Advertisement struct {
	ID uuid.UUID

	// metadata
	Name     string
	Topics   []string
	Coverage int
	Budget   int
	Message  string
	Status   models.AdvertisementStatus
}

type TopicWithCoverage struct {
	Name     string
	Coverage int
}
