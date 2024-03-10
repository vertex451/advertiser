package models

import (
	"gorm.io/gorm"
	"time"
)

type ChannelTopic struct {
	ChannelID int64  `gorm:"primaryKey"`
	TopicID   string `gorm:"primaryKey"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
