package postgresql

import (
	"advertiser/channel_owner/internal/config"
	"advertiser/shared/pkg/service/repo"
	"gorm.io/gorm"
	"sync"
)

type Repository struct {
	Db                *gorm.DB
	channelIdByHandle sync.Map // map[channelHandle]channelID
}

func New(cfg *config.Config) *Repository {
	return &Repository{
		Db: repo.New(cfg.PostgresHost),
	}
}
