package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type AdvertisementStatus string

const (
	// AdsStatusCreated - initial status when agency creates an advertisement
	AdsStatusCreated AdvertisementStatus = "created"
	// AdsStatusPending - status when agency clicks run button and triggers AdvertisementChannel creation
	AdsStatusPending AdvertisementStatus = "pending"
	// AdsStatusRunning - status when AdvertisementChannel entries are created
	AdsStatusRunning     AdvertisementStatus = "running"
	AdsStatusPaused      AdvertisementStatus = "paused"
	AdsStatusOutOfBudget AdvertisementStatus = "out_of_budget"
	AdsStatusFinished    AdvertisementStatus = "finished"
)

type Advertisement struct {
	ID         uuid.UUID `gorm:"primary_key"`
	CampaignID uuid.UUID

	// metadata
	Name         string
	TargetTopics []Topic `gorm:"many2many:advertisement_topics"`
	Message      string
	Status       AdvertisementStatus

	Budget      int
	CostPerView float32
	TotalViews  int
	TotalCost   float32

	// statistics
	AdsChannel []AdvertisementChannel `gorm:"foreignKey:AdvertisementID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (a *Advertisement) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.NewV4()
	}

	return
}
