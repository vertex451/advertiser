package postgresql

import (
	"advertiser/shared/config/config"
	"advertiser/shared/pkg/service/repo"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

func New(cfg *config.Config) *Repository {
	return &Repository{
		Db: repo.New(cfg),
	}
}
