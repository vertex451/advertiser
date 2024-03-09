package integration_tests

import (
	"advertiser/channel_owner/internal/config"
	"advertiser/channel_owner/internal/service/listener/repo/postgresql"
	"advertiser/channel_owner/internal/service/listener/transport"
	"advertiser/channel_owner/internal/service/listener/transport/mocks"
	"advertiser/channel_owner/internal/service/listener/usecase"
	"advertiser/shared/pkg/logger"
	"advertiser/shared/pkg/service/repo/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
)

func startChannelOwnerService(updatesChan chan tgbotapi.Update, targetChan chan tgbotapi.Chattable) {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		zap.L().Panic("error loading config", zap.Error(err))
	}

	logger.Init(cfg.LogLevel)

	tgBotApi := mocks.TgBotApiMock{
		UpdatesChan: updatesChan,
		TargetChan:  targetChan,
	}

	r := postgresql.New(cfg)

	uc := usecase.New(r)
	t := transport.New(uc, tgBotApi)

	go t.MonitorChannels()
	go t.RunNotificationService()

	fillTopics(r.Db)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
}

func fillTopics(db *gorm.DB) {
	topics := []models.Topic{{ID: "art"}, {ID: "books"}, {ID: "food"}, {ID: "pets"}, {ID: "sport"}}

	for _, topic := range topics {
		result := db.FirstOrCreate(&topic)
		if result.Error != nil {
			zap.L().Panic("failed to create topic", zap.Error(result.Error))
		}
	}
}
