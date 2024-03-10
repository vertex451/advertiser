package integration_tests

import (
	"advertiser/channel_owner/internal/config"
	"advertiser/channel_owner/internal/service/listener/repo/postgresql"
	"advertiser/channel_owner/internal/service/listener/transport"
	"advertiser/channel_owner/internal/service/listener/transport/mocks"
	"advertiser/channel_owner/internal/service/listener/usecase"
	"advertiser/shared/pkg/logger"
	"advertiser/shared/pkg/service/repo"
	"advertiser/shared/pkg/service/repo/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	// we get topics from UseCase cache, so don't move it after it.
	fillTopics(r.Db)

	uc := usecase.New(r)
	t := transport.New(uc, tgBotApi, cfg.Env)

	go t.MonitorChannels()
	go t.RunNotificationService()
}

func deleteTables(db *gorm.DB) {
	fmt.Println("### deleteTables!")
	var err error
	for _, table := range repo.GetAllTables() {
		err = db.Migrator().DropTable(table)
		if err != nil {
			zap.L().Panic("failed to drop table", zap.Error(err))
		}
	}
}

func fillTopics(db *gorm.DB) {
	topics := []models.Topic{{ID: "art"}, {ID: "books"}, {ID: "food"}, {ID: "pets"}, {ID: "sport"}}
	for _, topic := range topics {
		zap.L().Info("filling topic", zap.String("topic", topic.ID))
		result := db.FirstOrCreate(&topic)
		if result.Error != nil {
			zap.L().Panic("failed to create topic", zap.Error(result.Error))
		}
	}
}
