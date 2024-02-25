package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type AdvertisementStatus string

const (
	AdsStatusCreated     AdvertisementStatus = "created"
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
	Budget       int
	Message      string
	Status       AdvertisementStatus

	// statistics

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (a *Advertisement) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewV4()

	return
}
