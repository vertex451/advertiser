package dep_container

import (
	"advertiser/channel_owner/internal/config"
	"advertiser/channel_owner/internal/service/listener/repo/postgresql"
	"advertiser/channel_owner/internal/service/listener/transport"
	"advertiser/channel_owner/internal/service/listener/usecase"
	"advertiser/shared/tg_bot_api"
	"github.com/sarulabs/di"
)

const listenerServiceDefName = "listener-service"

// RegisterListenerService registers RegisterListenerService dependency.
func RegisterListenerService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: listenerServiceDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(configDefName).(*config.Config)
			tgBotApi := tg_bot_api.New(cfg.TelegramToken)
			r := ctn.Get(postgresqlDefName).(*postgresql.Repository)
			uc := usecase.New(r)

			return transport.New(uc, tgBotApi), nil
		},
	})
}

// MonitorChannels runs MonitorChannels dependency.
func (c Container) MonitorChannels() {
	c.container.Get(listenerServiceDefName).(*transport.Service).MonitorChannels()
}

// RunNotificationService runs RunNotificationService dependency.
func (c Container) RunNotificationService() {
	c.container.Get(listenerServiceDefName).(*transport.Service).RunNotificationService()
}
