package dep_container

import (
	"advertiser/shared/config/config"
	"github.com/sarulabs/di"
	"tg-bot/internal/service/bot_api/repo/postgresql"
)

const postgresqlDefName = "postgresql"

// RegisterPostgresqlService ...
func RegisterPostgresqlService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: postgresqlDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(configDefName).(*config.Config)
			return postgresql.New(cfg), nil
		},
	})
}
