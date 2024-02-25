package dep_container

import (
	"advertiser/channel_owner/internal/service/listener/repo/postgresql"
	"advertiser/channel_owner/internal/service/listener/usecase"
	"advertiser/channel_owner/internal/service/notification/transport"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sarulabs/di"
)

const notificationServiceDefName = "notification-service"

// RegisterNotificationService registers RegisterNotificationService dependency.
func RegisterNotificationService(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: notificationServiceDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			tgBotApi := ctn.Get(tbBotApiDefName).(*tgbotapi.BotAPI)

			r := ctn.Get(postgresqlDefName).(*postgresql.Repository)
			uc := usecase.New(r)
			return transport.New(uc, tgBotApi), nil
		},
	})
}
