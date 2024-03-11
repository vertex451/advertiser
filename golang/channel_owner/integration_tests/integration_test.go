package integration_tests

import (
	"advertiser/shared/pkg/service/repo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
			suite.updatesChan <- *tt.update
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
			suite.updatesChan <- *tt.update
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
			suite.updatesChan <- *tt.update
			msgRaw, _ := <-suite.targetChan
			msg := msgRaw.(tgbotapi.MessageConfig)
			assert.Equal(suite.T(), tt.expectedMsgText, msg.Text)
		})
	}
}

func (suite *MyTestSuite) TestModerationAndAdPosting() {
	suite.PrepareModerationTest()

	tc := []testCase{
		{
			testName: "TestModerateNotification",
			update:   nil,
			expectedMsgText: `
You have 1 advertisements to moderate.
Click on /moderate to view them.
`,
		},
		{
			testName:        "TestModerateCommand",
			update:          moderateCommandUpdate(),
			expectedMsgText: "Select an advertisement to moderate:",
		},
		{
			testName:        "TestModerateCallback",
			update:          moderateCallbackUpdate(),
			expectedMsgText: "Select an advertisement to moderate:",
		},
	}

	for _, tt := range tc {
		suite.Run(tt.testName, func() {
			if tt.update != nil {
				suite.updatesChan <- *tt.update
			}
			msgRaw, _ := <-suite.targetChan
			msg := msgRaw.(tgbotapi.MessageConfig)
			assert.Equal(suite.T(), tt.expectedMsgText, msg.Text)
		})
	}
}
