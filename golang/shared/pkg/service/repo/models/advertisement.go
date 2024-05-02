package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type AdvertisementStatus string

const (
	// AdsStatusCreated - initial status when agency creates an advertisement
	AdsStatusCreated AdvertisementStatus = "created"
	// AdsStatusPending - status when agency clicks run button and triggers AdvertisementChannel creation
	AdsStatusPending AdvertisementStatus = "pending"
	// AdsStatusRunning - status when AdvertisementChannel entries are created
	AdsStatusRunning     AdvertisementStatus = "running"
	AdsStatusPaused      AdvertisementStatus = "paused"
	AdsStatusOutOfBudget AdvertisementStatus = "out_of_budget"
	AdsStatusFinished    AdvertisementStatus = "finished"
)

type Advertisement struct {
	ID         uuid.UUID `gorm:"primary_key"`
	CampaignID uuid.UUID

	// metadata
	Name         string
	TargetTopics []Topic `gorm:"many2many:advertisement_topics"`

	Status AdvertisementStatus

	Budget      int
	CostPerMile float32
	TotalViews  int
	TotalCost   float32

	// Message
	MsgText     string
	MsgEntities []MsgEntity
	MsgImageURL string

	// statistics
	AdsChannel []AdvertisementChannel `gorm:"foreignKey:AdvertisementID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type MsgEntity struct {
	ID              uuid.UUID `gorm:"primary_key"`
	AdvertisementID uuid.UUID

	Type     string
	Offset   int
	Length   int
	URL      string
	Language string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (a *Advertisement) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.NewV4()
	}

	return
}

func (a *Advertisement) GetTopics() string {
	var topics []string
	for _, topic := range a.TargetTopics {
		topics = append(topics, topic.ID)
	}

	return strings.Join(topics, ", ")
}

func (e *MsgEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.NewV4()
	}

	return
}
