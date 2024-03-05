package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"time"
)

type AdChanStatus string

const (
	AdChanCreated           AdChanStatus = "created"
	AdChanWaitingToBePosted AdChanStatus = "waiting_to_be_posted"
	AdChanPosted            AdChanStatus = "posted"
	AdChanRejected          AdChanStatus = "rejected"
	AdChanFinished          AdChanStatus = "finished"
)

type AdvertisementChannel struct {
	ID              uuid.UUID `gorm:"primaryKey"`
	AdvertisementID uuid.UUID `gorm:"index:idx_advertisement_channel,unique"`
	ChannelID       int64     `gorm:"index:idx_advertisement_channel,unique"`

	Status          AdChanStatus
	RejectionReason string

	ChannelTitle   string
	ChannelHandle  string
	ChannelOwnerID int64

	AdName        string
	AdMessage     string
	AdCostPerView float32

	// stats
	MessageID           int
	DeletionScheduledAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ac *AdvertisementChannel) BeforeCreate(tx *gorm.DB) (err error) {
	if ac.ID == uuid.Nil {
		ac.ID = uuid.NewV4()
	}

	return
}
