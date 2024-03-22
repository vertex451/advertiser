package transport

//
//func TestAllTopicsCallback(t *testing.T) {
//	updatesChan := make(chan tgbotapi.Update)
//	targetChan := make(chan tgbotapi.Chattable)
//
//	s := New(UseCaseMock{}, TgBotApiMock{UpdatesChan: updatesChan, TargetChan: targetChan})
//	go s.MonitorChannels()
//
//	tc := []struct {
//		testName        string
//		update          tgbotapi.Update
//		expectedMsgText string
//	}{
//		{
//			testName: "TestAllTopicsCommand",
//			update: tgbotapi.Update{
//				Message: &tgbotapi.Message{
//					MsgEntities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: len(constants.AllTopics) + 1}},
//					Chat:     &tgbotapi.Chat{ID: 6406834985},
//					MsgText:     fmt.Sprintf("/%s", constants.AllTopics)},
//			},
//			expectedMsgText: `
//Supported topics:
//topic1, topic2, topic3
//`,
//		},
//		{
//			testName: "TestAllTopicsCallback",
//			update: tgbotapi.Update{
//				CallbackQuery: &tgbotapi.CallbackQuery{
//					Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 6406834985}}, Data: constants.AllTopics,
//				},
//			},
//			expectedMsgText: `
//Supported topics:
//topic1, topic2, topic3
//`,
//		},
//		{
//			testName: "TestStartCommand",
//			update: tgbotapi.Update{
//				Message: &tgbotapi.Message{
//					MsgEntities: []tgbotapi.MessageEntity{{Type: "bot_command", Length: len(constants.Start) + 1}},
//					Chat:     &tgbotapi.Chat{ID: 6406834985},
//					MsgText:     fmt.Sprintf("/%s", constants.Start)},
//			},
//			expectedMsgText: "Choose action:",
//		},
//		{
//			testName: "TestStartCallback",
//			update: tgbotapi.Update{
//				CallbackQuery: &tgbotapi.CallbackQuery{
//					Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 6406834985}}, Data: constants.Start,
//				},
//			},
//			expectedMsgText: "Choose action:",
//		},
//		{
//			testName: "TestBackCallback",
//			update: tgbotapi.Update{
//				CallbackQuery: &tgbotapi.CallbackQuery{
//					Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 6406834985}}, Data: constants.Back,
//				},
//			},
//		},
//	}
//
//	for _, tt := range tc {
//		t.Run(tt.testName, func(t *testing.T) {
//			updatesChan <- tt.update
//			msgRaw, _ := <-targetChan
//			msg := msgRaw.(tgbotapi.MessageConfig)
//			assert.Equal(t, tt.expectedMsgText, msg.MsgText)
//		})
//	}
//}
