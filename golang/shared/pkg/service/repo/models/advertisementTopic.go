package models

import uuid "github.com/satori/go.uuid"

type AdvertisementTopic struct {
	AdvertisementID uuid.UUID `gorm:"primaryKey"`
	TopicID         string    `gorm:"primaryKey"`
}
