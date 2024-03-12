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

func (suite *MyTestSuite) TestCreateCampaign() {
	tc := []testCase{
		{
			testName:        "TestCreateCampaignCallback",
			update:          createCampaignCallbackUpdate,
			expectedMsgText: "Enter campaign name:",
		},
		{ // state dependant test
			testName:        "TestCreateCampaignMessage",
			update:          createCampaignMessageUpdate,
			expectedMsgText: "Campaign Food created!",
		},
		{ // state dependant test
			testName:        "TestMyCampaigns",
			update:          myCampaignsCallbackUpdate,
			expectedMsgText: "Select a campaign:",
		},
		{ // state dependant test
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
		{ // state dependant test
			testName:        "TestCreateAdvertisementMessage",
			update:          createAdMessageUpdate,
			expectedMsgText: "Advertisement  created!",
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
