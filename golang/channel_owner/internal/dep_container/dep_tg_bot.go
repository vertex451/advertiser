package dep_container

import (
	"advertiser/channel_owner/internal/config"
	"advertiser/channel_owner/internal/service/repo/postgresql"
	"advertiser/channel_owner/internal/service/transport"
	"advertiser/channel_owner/internal/service/usecase"
	"github.com/sarulabs/di"
)

const tbBotApiName = "tg-bot-api"

// RegisterTgBotApiService registers RegisterTgBotApiService dependency.
func RegisterTgBotApiService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: tbBotApiName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(configDefName).(*config.Config)
			r := ctn.Get(postgresqlDefName).(*postgresql.Repository)
			uc := usecase.New(r)
			return transport.New(uc, cfg.TelegramToken), nil
		},
	})
}

// RunChannelListener runs RunChannelListener dependency.
func (c Container) RunChannelListener() {
	c.container.Get(tbBotApiName).(*transport.Transport).MonitorChannels()
}
