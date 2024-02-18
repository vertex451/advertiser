package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Advertisement struct {
	ID         uuid.UUID `gorm:"primary_key"`
	CampaignID uuid.UUID

	// metadata
	Name         string
	TargetTopics []Topic `gorm:"many2many:advertisement_topics"`
	Budget       int
	Message      string

	// statistics

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (a *Advertisement) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.NewV4()

	return
}
