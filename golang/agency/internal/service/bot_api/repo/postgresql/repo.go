package postgresql

import (
	"advertiser/shared/pkg/service/repo"
	"gorm.io/gorm"
)

type Repository struct {
	Db *gorm.DB
}

func New() *Repository {
	return &Repository{
		Db: repo.New(),
	}
}
