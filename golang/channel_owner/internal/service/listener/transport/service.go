package transport

import (
	"advertiser/channel_owner/internal/service/listener"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"sync"
)

type Transport struct {
	uc           listener.UseCase
	tgBotApi     *tgbotapi.BotAPI
	updateConfig tgbotapi.UpdateConfig
	state        sync.Map // map[ChatID]stateData
	lastMsgID    int
}

func New(uc listener.UseCase, tgBotApi *tgbotapi.BotAPI) *Transport {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	return &Transport{
		tgBotApi:     tgBotApi,
		updateConfig: updateConfig,
		uc:           uc,
	}
}

func (s *Transport) MonitorChannels() {
	updates := s.tgBotApi.GetUpdatesChan(s.updateConfig)

	for update := range updates {

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the bot! Press the button below to get started.")
		//msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		//	tgbotapi.NewKeyboardButtonRow(
		//		tgbotapi.NewKeyboardButton("start"),
		//	))
		//sentMsg, err := s.tgBotApi.Send(msg)
		//if err != nil {
		//	zap.L().Error("failed to send message", zap.Error(err))
		//}

		var err error
		var sentMsg tgbotapi.Message
		responseMessage := s.handleUpdate(update)
		if responseMessage != nil {
			if s.lastMsgID != 0 {
				deleteMsg := tgbotapi.NewDeleteMessage(TgBotDirectChatID, s.lastMsgID)
				s.tgBotApi.Send(deleteMsg)
				s.lastMsgID = 0
			}
			sentMsg, err = s.tgBotApi.Send(responseMessage)
			if err != nil {
				zap.L().Error("failed to send message", zap.Error(err))
				continue
			}
			s.lastMsgID = sentMsg.MessageID
		}
	}
}

// TODO add rate limiter
func (s *Transport) handleUpdate(update tgbotapi.Update) *tgbotapi.MessageConfig {
	if update.Message != nil {
		if update.Message.IsCommand() {
			return s.handleCommand(update)
		}

		if update.Message.Text != "" {
			return s.handleStateQuery(update)
		}
	}

	if update.CallbackQuery != nil {
		return s.handleCallbackQuery(update.CallbackQuery)
	}

	if update.MyChatMember != nil {
		return s.handleUpdateEvent(update)
	}

	return nil
}

// Added bot by admin to channel
// {"ok":true,"result":[{"update_id":632156492, "my_chat_member":{"chat":{"id":-1002025237232,"title":"Public Sport Channel","username":"sportchannel451","type":"channel"},"from":{"id":399749369,"is_bot":false,"first_name":"Artem","username":"vertex451","language_code":"en"},"date":1705664974,"old_chat_member":{"user":{"id":6845534569,"is_bot":true,"first_name":"Advertiser","username":"advertiser_451_bot"},"status":"left"},"new_chat_member":{"user":{"id":6845534569,"is_bot":true,"first_name":"Advertiser","username":"advertiser_451_bot"},"status":"administrator","can_be_edited":false,"can_manage_chat":true,"can_change_info":false,"can_post_messages":true,"can_edit_messages":false,"can_delete_messages":false,"can_invite_users":false,"can_restrict_members":true,"can_promote_members":false,"can_manage_video_chats":false,"can_post_stories":false,"can_edit_stories":false,"can_delete_stories":false,"is_anonymous":false,"can_manage_voice_chats":false}}}]}
// Remove bot from admins
// {"ok":true,"result":[{"update_id":632156491, "my_chat_member":{"chat":{"id":-1002025237232,"title":"Public Sport Channel","username":"sportchannel451","type":"channel"},"from":{"id":399749369,"is_bot":false,"first_name":"Artem","username":"vertex451","language_code":"en"},"date":1705664660,"old_chat_member":{"user":{"id":6845534569,"is_bot":true,"first_name":"Advertiser","username":"advertiser_451_bot"},"status":"administrator","can_be_edited":false,"can_manage_chat":true,"can_change_info":false,"can_post_messages":true,"can_edit_messages":false,"can_delete_messages":false,"can_invite_users":false,"can_restrict_members":true,"can_promote_members":false,"can_manage_video_chats":false,"can_post_stories":false,"can_edit_stories":false,"can_delete_stories":false,"is_anonymous":false,"can_manage_voice_chats":false},"new_chat_member":{"user":{"id":6845534569,"is_bot":true,"first_name":"Advertiser","username":"advertiser_451_bot"},"status":"left"}}}]}
