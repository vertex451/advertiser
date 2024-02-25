package dep_container

import (
	"advertiser/channel_owner/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sarulabs/di"
)

const tbBotApiDefName = "tg-bot-api"

// RegisterTgBotApi registers RegisterTgBotApi dependency.
func RegisterTgBotApi(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: tbBotApiDefName,
		Build: func(ctn di.Container) (interface{}, error) {
			cfg := ctn.Get(configDefName).(*config.Config)
			tgBotApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
			if err != nil {
				panic(err)
			}
			tgBotApi.Debug = true

			return tgBotApi, nil
		},
	})
}
