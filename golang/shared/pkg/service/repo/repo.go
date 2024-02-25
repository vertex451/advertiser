package repo

import (
	models2 "advertiser/shared/pkg/service/repo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Warsaw"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		models2.Channel{},
		models2.Topic{},
		models2.User{},
		models2.ChannelAdmin{},
		models2.Campaign{},
		models2.Advertisement{},
		models2.AdvertisementChannel{},
	)
	if err != nil {
		panic("failed to AutoMigrate")
	}

	return db
}
