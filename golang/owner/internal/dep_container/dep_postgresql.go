package dep_container

import (
	"advertiser/owner/internal/service/listener/repo/postgresql"
	"advertiser/shared/config/config"
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
