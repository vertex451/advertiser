package repo

import (
	"advertiser/shared/pkg/service/repo/models"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetAllTables() []interface{} {
	return []interface{}{
		models.Topic{},
		models.User{},
		models.Channel{},
		models.ChannelTopic{},
		models.ChannelAdmin{},
		models.Campaign{},
		models.Advertisement{},
		models.AdvertisementChannel{},
		models.AdvertisementTopic{},
	}
}

func New(host string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Warsaw",
		host,
		"postgres",
		"postgres",
		"postgres",
		"5432",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Panic("failed to connect database", zap.Error(err))
	}

	err = db.AutoMigrate(
		GetAllTables()...,
	)
	if err != nil {
		zap.L().Panic("failed to AutoMigrate", zap.Error(err))
	}

	return db
}
