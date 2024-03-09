package models

import (
	"gorm.io/gorm"
	"time"
)

type Channel struct {
	ID int64 `gorm:"primary_key"`

	//Metadata
	Admins      []*User  `gorm:"many2many:channel_admins;"`
	Topics      []*Topic `gorm:"many2many:channel_topics;"`
	Description string
	Handle      string `gorm:"index"`

	IsChannel      bool
	RewardsAddress string
	Title          string

	//Stats
	PostsPerDay               int
	Subscribers               int
	TotalViewsForLastTenPosts int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
