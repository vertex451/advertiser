package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type AdvertisementTopic struct {
	AdvertisementID uuid.UUID `gorm:"primaryKey"`
	TopicID         string    `gorm:"primaryKey"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
