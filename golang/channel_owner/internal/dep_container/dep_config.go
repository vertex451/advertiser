package dep_container

import (
	"advertiser/shared/config/config"
	"github.com/sarulabs/di"
)

const configDefName = "config"

// RegisterConfig registers Config dependency.
func RegisterConfig(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: configDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			return config.Load(), nil
		},
	})
}
