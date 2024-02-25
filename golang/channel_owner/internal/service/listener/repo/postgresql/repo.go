package postgresql

import (
	"advertiser/shared/pkg/service/repo"
	"gorm.io/gorm"
	"sync"
)

type Repository struct {
	Db                *gorm.DB
	channelIdByHandle sync.Map // map[channelHandle]channelID
}

func New() *Repository {
	return &Repository{
		Db: repo.New(),
	}
}
