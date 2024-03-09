package models

import (
	"gorm.io/gorm"
	"time"
)

type ChannelAdmin struct {
	ChannelID int64 `gorm:"primaryKey"`
	UserID    int64 `gorm:"primaryKey"`
	Role      string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
