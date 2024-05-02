package models

import (
	"gorm.io/gorm"
	"time"
)

type Topic struct {
	ID string `gorm:"primary_key;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

//func (t Topic) Join() string {
//	return "topics"
//}
