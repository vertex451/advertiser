package integration_tests

import (
	"advertiser/channel_owner/internal/config"
	"advertiser/channel_owner/internal/service/listener/repo/postgresql"
	"advertiser/channel_owner/internal/service/listener/transport/mocks"
	"advertiser/shared/pkg/service/repo"
	"advertiser/shared/pkg/service/repo/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type MyTestSuite struct {
	suite.Suite
	updatesChan chan tgbotapi.Update
	targetChan  chan tgbotapi.Chattable
}

// SetupTest is called before each test in the suite
func (suite *MyTestSuite) SetupTest() {
	// Add setup logic here
	suite.T().Log("SetupTest: This runs before each test")

	suite.updatesChan = make(chan tgbotapi.Update)
	suite.targetChan = make(chan tgbotapi.Chattable)

	startChannelOwnerService(suite.updatesChan, suite.targetChan)
}

// TearDownTest is called after each test in the suite
func (suite *MyTestSuite) TearDownTest() {
	// Add teardown logic here

	db := repo.New("localhost")
	defer deleteTables(db)

	//close(suite.updatesChan)
	//close(suite.targetChan)

	suite.T().Log("TearDownTest: This runs after each test")
}

func TestMyTestSuite(t *testing.T) {
	// Run the test suite
	suite.Run(t, new(MyTestSuite))
}

func (suite *MyTestSuite) TestStart() {
	tc := []testCase{
		{
			testName:        "TestStartCommand",
			update:          startCommandUpdate(),
			expectedMsgText: "Choose action:",
		},
		{
			testName:        "TestStartCallback",
			update:          startCallbackUpdate(),
			expectedMsgText: "Choose action:",
		},
	}

	for _, tt := range tc {
		suite.Run(tt.testName, func() {
			suite.updatesChan <- tt.update
			msgRaw, _ := <-suite.targetChan
			msg := msgRaw.(tgbotapi.MessageConfig)
			assert.Equal(suite.T(), tt.expectedMsgText, msg.Text)
		})
	}
}

func (suite *MyTestSuite) TestAllTopics() {

	tc := []testCase{
		{
			testName: "TestAllTopicsCommand",
			update:   allTopicsCommandUpdate(),
			expectedMsgText: `
	Supported topics:
	art, books, food, pets, sport
	`,
		},
		{
			testName: "TestAllTopicsCallback",
			update:   allTopicsCallbackUpdate(),
			expectedMsgText: `
	Supported topics:
	art, books, food, pets, sport
	`,
		},
	}

	for _, tt := range tc {
		suite.Run(tt.testName, func() {
			suite.updatesChan <- allTopicsCallbackUpdate()
			msgRaw, _ := <-suite.targetChan
			msg := msgRaw.(tgbotapi.MessageConfig)
			assert.Equal(suite.T(),
				`
Supported topics:
art, books, food, pets, sport
`, msg.Text)
		})
	}
}

// TestChannelOwnerFlow
// 1. Handle event of adding bot to channel
// 2. Set topics
// 3. List my channels
func (suite *MyTestSuite) TestChannelOwnerFlow() {
	tc := []testCase{
		{
			testName:        "TestBotIsAddedToChannelEvent",
			update:          botIsAddedToChannelUpdate(),
			expectedMsgText: "Advertiser bot was successfully added to Sport Channel",
		},
		{
			testName:        "TestMyChannelsCallback",
			update:          myChannelsCallbackUpdate(),
			expectedMsgText: "Select a channel:",
		},
		{
			testName: "TestEditTopicsCallback",
			update:   editTopicsCallbackUpdate(),
			expectedMsgText: `
Choose topics from the list:
art, books, food, pets, sport
`,
		},
		{
			// TestEditTopicsMessage must go after TestEditTopicsCallback due to state nature
			testName:        "TestEditTopicsMessage",
			update:          editTopicsMessageUpdate(),
			expectedMsgText: "Topics changed! New channel topics: art, books, food, pets, sport",
		},
		{
			testName:        "TestBotIsRemovedFromChannelEvent",
			update:          botIsRemovedFromChannelUpdate(),
			expectedMsgText: "Advertiser bot is removed from Sport Channel",
		},
	}

	for _, tt := range tc {
		suite.Run(tt.testName, func() {
			suite.updatesChan <- tt.update
			msgRaw, _ := <-suite.targetChan
			msg := msgRaw.(tgbotapi.MessageConfig)
			assert.Equal(suite.T(), tt.expectedMsgText, msg.Text)
		})
	}
}

func (suite *MyTestSuite) TestModerationAndAdPosting() {
	suite.updatesChan <- botIsAddedToChannelUpdate()
	<-suite.targetChan

	suite.updatesChan <- editTopicsCallbackUpdate()
	<-suite.targetChan
	suite.updatesChan <- editTopicsMessageUpdate()
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

	msgRaw, _ := <-suite.targetChan
	msg := msgRaw.(tgbotapi.MessageConfig)
	assert.Equal(suite.T(), `
You have 1 advertisements to moderate.
Click on /moderate to view them.
`, msg.Text)

}
