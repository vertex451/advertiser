package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type FeatureRequest struct {
	ID uuid.UUID `gorm:"primaryKey"`

	RequestedBy int64
	Message     string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (fr *FeatureRequest) BeforeCreate(tx *gorm.DB) (err error) {
	if fr.ID == uuid.Nil {
		fr.ID = uuid.NewV4()
	}

	return
}
