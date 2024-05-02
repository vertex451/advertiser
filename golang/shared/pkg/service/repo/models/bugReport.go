package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type BugReport struct {
	ID uuid.UUID `gorm:"primaryKey"`

	ReportedBy int64
	Message    string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (bg *BugReport) BeforeCreate(tx *gorm.DB) (err error) {
	if bg.ID == uuid.Nil {
		bg.ID = uuid.NewV4()
	}

	return
}
