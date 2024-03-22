package dep_container

import (
	"advertiser/owner/internal/service/listener/repo/postgresql"
	"advertiser/owner/internal/service/listener/transport"
	"advertiser/owner/internal/service/listener/usecase"
	"advertiser/shared/config/config"
	"advertiser/shared/tg_bot_api"
	"github.com/sarulabs/di"
)

const ownerServiceDefName = "owner-service"

// RegisterListenerService registers RegisterListenerService dependency.
func RegisterListenerService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: ownerServiceDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(configDefName).(*config.Config)
			tgBotApi := tg_bot_api.New(cfg.Secrets.OwnerTgToken)
			r := ctn.Get(postgresqlDefName).(*postgresql.Repository)
			uc := usecase.New(r)

			return transport.New(uc, tgBotApi, cfg.Env), nil
		},
	})
}

// MonitorChannels runs MonitorChannels dependency.
func (c Container) MonitorChannels() {
	c.container.Get(ownerServiceDefName).(*transport.Service).MonitorChannels()
}

// RunNotificationService runs RunNotificationService dependency.
func (c Container) RunNotificationService() {
	c.container.Get(ownerServiceDefName).(*transport.Service).RunNotificationService()
}
