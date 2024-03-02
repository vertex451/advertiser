package dep_container

import (
	"advertiser/channel_owner/internal/service/listener/repo/postgresql"
	"advertiser/channel_owner/internal/service/listener/transport"
	"advertiser/channel_owner/internal/service/listener/usecase"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sarulabs/di"
)

const listenerServiceDefName = "listener-service"

// RegisterListenerService registers RegisterListenerService dependency.
func RegisterListenerService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: listenerServiceDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			tgBotApi := ctn.Get(tbBotApiDefName).(*tgbotapi.BotAPI)
			r := ctn.Get(postgresqlDefName).(*postgresql.Repository)
			uc := usecase.New(r)

			return transport.New(uc, tgBotApi), nil
		},
	})
}

// RunChannelListener runs RunChannelListener dependency.
func (c Container) RunChannelListener() {
	c.container.Get(listenerServiceDefName).(*transport.Transport).MonitorChannels()
}

// RunNotificationService runs RunNotificationService dependency.
func (c Container) RunNotificationService() {
	c.container.Get(listenerServiceDefName).(*transport.Transport).RunNotificationService()
}
