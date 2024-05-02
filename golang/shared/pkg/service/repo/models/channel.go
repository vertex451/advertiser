package models

import (
	"advertiser/shared/pkg/service/constants"
	"gorm.io/gorm"
	"time"
)

type Channel struct {
	ID int64 `gorm:"primary_key"`

	//Metadata
	Topics      []Topic `gorm:"many2many:channel_topics;"`
	CostPerMile float64
	Location    constants.Location

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
