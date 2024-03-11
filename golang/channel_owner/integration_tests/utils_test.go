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

func (suite *MyTestSuite) PrepareModerationTest() {
	suite.updatesChan <- *botIsAddedToChannelUpdate()
	<-suite.targetChan

	suite.updatesChan <- *editTopicsCallbackUpdate()
	<-suite.targetChan
	suite.updatesChan <- *editTopicsMessageUpdate()
	<-suite.targetChan

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		zap.L().Panic("error loading config", zap.Error(err))
	}
	r := postgresql.New(cfg)

	testCampaign := models.Campaign{
		UserID: mocks.ChannelCreator.ID,
		Name:   "TestCampaign",
	}

	err = r.Db.Create(&testCampaign).Error
	if err != nil {
		zap.L().Panic("failed to create test campaign", zap.Error(err))
	}
	testAd := models.Advertisement{
		CampaignID:   testCampaign.ID,
		Name:         "testAd",
		TargetTopics: []models.Topic{{ID: "food"}},
		Message:      "Visit our restaurant!",
		Status:       models.AdsStatusPending,
		Budget:       100,
		CostPerView:  0.01,
	}
	err = r.Db.Create(&testAd).Error
	if err != nil {
		zap.L().Panic("failed to create test ad", zap.Error(err))
	}
}
