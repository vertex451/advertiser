package repo

import (
	"advertiser/shared/pkg/service/repo/models"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	err = db.AutoMigrate(
		models.Channel{},
		models.Topic{},
		models.User{},
		models.ChannelAdmin{},
		models.Campaign{},
		models.Advertisement{},
		models.AdvertisementChannel{},
	)
	if err != nil {
		zap.L().Panic("failed to AutoMigrate", zap.Error(err))
	}

	return db
}
