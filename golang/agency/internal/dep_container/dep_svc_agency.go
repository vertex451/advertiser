package dep_container

import (
	"advertiser/shared/config/config"
	"advertiser/shared/pkg/storage/localstorage"
	"advertiser/shared/tg_bot_api"
	"github.com/sarulabs/di"
	"tg-bot/internal/service/bot_api/repo/postgresql"
	"tg-bot/internal/service/bot_api/transport"
	"tg-bot/internal/service/bot_api/usecase"
)

const tbBotApiName = "tg-bot-api"

// RegisterTgBotApiService registers RegisterTgBotApiService dependency.
func RegisterTgBotApiService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: tbBotApiName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(configDefName).(*config.Config)
			tgBotApi := tg_bot_api.New(cfg.Secrets.AgencyTgToken)
			r := ctn.Get(postgresqlDefName).(*postgresql.Repository)
			uc := usecase.New(r)
			return transport.New(uc, tgBotApi, localstorage.New(cfg.DataDir), cfg.Env), nil
		},
	})
}

// RunChannelListener runs RunChannelListener dependency.
func (c Container) RunChannelListener() {
	c.container.Get(tbBotApiName).(*transport.Transport).MonitorChannels()
}