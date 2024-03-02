package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type Campaign struct {
	ID             uuid.UUID `gorm:"primary_key"`
	UserID         int64
	Advertisements []Advertisement `gorm:"foreignKey:CampaignID"`

	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (c *Campaign) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.NewV4()
	}

	return
}
