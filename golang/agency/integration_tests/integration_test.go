package integration_tests

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func (suite *MyTestSuite) TestStart() {
	tc := []testCase{
		{
			testName:        "TestStartCommand",
			update:          startCommandUpdate,
			expectedMsgText: "Choose action:",
		},
		{
			testName:        "TestStartCallback",
			update:          startCallbackUpdate,
			expectedMsgText: "Choose action:",
		},
	}

	for _, tt := range tc {
		suite.Run(tt.testName, func() {
			suite.updatesChan <- *tt.update("")
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
			update:   allTopicsCommandUpdate,
			expectedMsgText: `
Supported topics:
art, books, food, pets, sport
`,
		},
		{
			testName: "TestAllTopicsWithCoverageCommand",
			update:   allTopicsWithCoverageCallbackUpdate,
			expectedMsgText: `art: 0 subscribers
books: 0 subscribers
food: 0 subscribers
pets: 0 subscribers
sport: 0 subscribers`,
		},
	}

	for _, tt := range tc {
		suite.Run(tt.testName, func() {
			suite.updatesChan <- *tt.update("")
			msgRaw, _ := <-suite.targetChan
			msg := msgRaw.(tgbotapi.MessageConfig)
			assert.Equal(suite.T(), tt.expectedMsgText, msg.Text)
		})
	}
}

// TestCreateCampaign is state dependant test, single TC may fail if run separately
func (suite *MyTestSuite) TestCreateCampaign() {
	tc := []testCase{
		{
			testName:        "TestCreateCampaignCallback",
			update:          createCampaignCallbackUpdate,
			expectedMsgText: "Enter campaign name:",
		},
		{
			testName:        "TestCreateCampaignMessage",
			update:          createCampaignMessageUpdate,
			expectedMsgText: "Campaign Food created!",
		},
		{
			testName:        "TestMyCampaigns",
			update:          myCampaignsCallbackUpdate,
			expectedMsgText: "Select a campaign:",
		},
		{
			testName: "TestCreateAdvertisement",
			preHook:  getCampaignID,
			update:   createAdCallbackUpdate,
			expectedMsgText: `
To create an advertisement, send a message in the following format:
Name: Stock market
TargetTopics: topic1, topic2, topic3
BudgetUSD: 100
CostPerView: 0.1
Message: Follow this [link](https://www.investing.com/) to find more about investments!
`,
		},
		{
			testName:        "TestCreateAdvertisementMessage",
			update:          createAdMessageUpdate,
			expectedMsgText: "Advertisement 'Stock market' created!",
		},
		{
			testName:        "TestCampaignDetails",
			preHook:         getCampaignID,
			update:          campaignDetailsCallbackUpdate,
			expectedMsgText: "Food advertisements:",
		},
		{
			testName: "TestAdDetails",
			preHook:  getAdID,
			update:   adDetailsCallbackUpdate,
			expectedMsgText: `
Name: Stock market
Status: created
TargetTopics: art, food
BudgetUSD: 100
Message: Follow this [link](https://www.investing.com/) to find more about investments!

Total members: 0
`,
		},
		{
			testName:        "TestRunAd",
			preHook:         getAdID,
			update:          runAdCallbackUpdate,
			expectedMsgText: "Advertising is running! It will start appearing in channels after an approval from channel owners",
		},
	}

	for _, tt := range tc {
		suite.Run(tt.testName, func() {
			if tt.preHook != nil {
				suite.updatesChan <- *tt.update(tt.preHook(suite.db))
			} else {
				suite.updatesChan <- *tt.update("")
			}
			msgRaw, _ := <-suite.targetChan
			msg := msgRaw.(tgbotapi.MessageConfig)
			assert.Equal(suite.T(), tt.expectedMsgText, msg.Text)
		})
	}
}
