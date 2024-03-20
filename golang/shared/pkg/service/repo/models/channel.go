package models

import (
	"gorm.io/gorm"
	"time"
)

type Channel struct {
	ID int64 `gorm:"primary_key"`

	//Metadata
	Topics []Topic `gorm:"many2many:channel_topics;"`

	ChannelAdmins []ChannelAdmin `gorm:"foreignKey:ChannelID"`

	Description string
	Handle      string

	IsChannel      bool
	RewardsAddress string
	Title          string

	//Stats
	Subscribers int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
