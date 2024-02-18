package dep_container

import (
	"github.com/sarulabs/di"
	"tg-bot/internal/service/bot_api/repo/postgresql"
)

const postgresqlDefName = "postgresql"

// RegisterPostgresqlService ...
func RegisterPostgresqlService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: postgresqlDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			return postgresql.New(), nil
		},
	})
}
