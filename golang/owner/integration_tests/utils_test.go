package integration_tests

import (
	"advertiser/owner/internal/service/listener/repo/postgresql"
	"advertiser/owner/internal/service/listener/transport"
	"advertiser/owner/internal/service/listener/usecase"
	"advertiser/shared/config/config"
	"advertiser/shared/pkg/logger"
	"advertiser/shared/pkg/service/repo"
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/tg_bot_api/mocks"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"testing"
)

type MyTestSuite struct {
	suite.Suite
	updatesChan chan tgbotapi.Update
	targetChan  chan tgbotapi.Chattable
	db          *gorm.DB
}

// SetupTest is called before each test in the suite
func (suite *MyTestSuite) SetupTest() {
	suite.startChannelOwnerService()
}

// TearDownTest is called after each test in the suite
func (suite *MyTestSuite) TearDownTest() {
	suite.deleteTables()

	//close(suite.updatesChan)
	//close(suite.targetChan)
}

func TestMyTestSuite(t *testing.T) {
	// Run the test suite
	suite.Run(t, new(MyTestSuite))
}

func (suite *MyTestSuite) InitMyTestSuite(cfg *config.Config) {
	suite.updatesChan = make(chan tgbotapi.Update)
	suite.targetChan = make(chan tgbotapi.Chattable)
	suite.db = repo.New(cfg)
}

func (suite *MyTestSuite) startChannelOwnerService() {
	cfg := config.Load()
	logger.Init(cfg.LogLevel)

	suite.InitMyTestSuite(cfg)

	tgBotApi := tgBotApiMock.TgBotApiMock{
		UpdatesChan: suite.updatesChan,
		TargetChan:  suite.targetChan,
	}

	r := postgresql.New(cfg)

	// we need it for tearDown function
	// we get topics from UseCase cache, so don't move fillTopics after usecase.New.
	fillTopics(r.Db)

	uc := usecase.New(r)
	t := transport.New(uc, tgBotApi, cfg.Env)

	go t.MonitorChannels()
	go t.RunNotificationService()
}

func (suite *MyTestSuite) deleteTables() {
	var err error
	for _, table := range repo.GetAllTables() {
		err = suite.db.Migrator().DropTable(table)
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

	r := postgresql.New(config.Load())

	testCampaign := models.Campaign{
		UserID: tgBotApiMock.ChannelCreator.ID,
		Name:   "TestCampaign",
	}

	err := r.Db.Create(&testCampaign).Error
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
