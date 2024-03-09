package dep_container

import (
	"advertiser/channel_owner/internal/config"
	"advertiser/channel_owner/internal/service/listener/repo/postgresql"
	"github.com/sarulabs/di"
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
