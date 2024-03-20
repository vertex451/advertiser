package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int64      `gorm:"primary_key"`
	Campaigns []Campaign `gorm:"foreignKey:user_id"`
	Channels  []Channel  `gorm:"many2many:channel_admins;"`

	Handle string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
