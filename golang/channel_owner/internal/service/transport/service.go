package transport

import (
	"advertiser/channel_owner/internal/service"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"sync"
)

const (
	TgBotDirectChatID int64 = 6406834985
)

// Commands:
const (
	Start              = "start"
	Back               = "back"
	AllTopics          = "all_topics"
	EditChannelsTopics = "edit_channel_topics"
	ListChannelsTopics = "list_channels_topics"
	MyChannels         = "my_channels"
)

type BotState int

const (
	StateStart BotState = iota
	StateEditTopics
	StateSetCampaignName
	StateCreateAd
	StateUpdateAd
)

type Transport struct {
	tgBotApi     *tgbotapi.BotAPI
	uc           service.UseCase
	updateConfig tgbotapi.UpdateConfig
	state        sync.Map // map[ChatID]stateData
	lastMsgID    int
}

func New(uc service.UseCase, tgToken string) *Transport {
	tgBotApi, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		panic(err)
	}
	tgBotApi.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	zap.L().Info("Started Transport")

	return &Transport{
		tgBotApi:     tgBotApi,
		uc:           uc,
		updateConfig: updateConfig,
	}
}

func (t *Transport) MonitorChannels() {
	updates := t.tgBotApi.GetUpdatesChan(t.updateConfig)

	for update := range updates {

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the bot! Press the button below to get started.")
		//msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		//	tgbotapi.NewKeyboardButtonRow(
		//		tgbotapi.NewKeyboardButton("Start"),
		//	))
		//sentMsg, err := t.tgBotApi.Send(msg)
		//if err != nil {
		//	zap.L().Error("failed to send message", zap.Error(err))
		//}

		var err error
		var sentMsg tgbotapi.Message
		responseMessage := t.handleUpdate(update)
		if responseMessage != nil {
			if t.lastMsgID != 0 {
				deleteMsg := tgbotapi.NewDeleteMessage(TgBotDirectChatID, t.lastMsgID)
				t.tgBotApi.Send(deleteMsg)
				t.lastMsgID = 0
			}
			sentMsg, err = t.tgBotApi.Send(responseMessage)
			if err != nil {
				zap.L().Error("failed to send message", zap.Error(err))
				continue
			}
			t.lastMsgID = sentMsg.MessageID
		}
	}
}

// TODO add rate limiter
func (t *Transport) handleUpdate(update tgbotapi.Update) *tgbotapi.MessageConfig {
	if update.Message != nil {
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case Start:
				return t.start(update.Message.Chat.ID)
			case AllTopics:
				return t.allTopics(update.Message.Chat.ID)
			}
		}

		if update.Message.Text != "" {
			return t.handleStateQuery(update)
		}
	}

	if update.CallbackQuery != nil {
		return t.handleCallbackQuery(update.CallbackQuery)
	}

	if update.MyChatMember != nil {
		return t.handleUpdateEvent(update)
	}

	return nil
}

// Added bot by admin to channel
// {"ok":true,"result":[{"update_id":632156492, "my_chat_member":{"chat":{"id":-1002025237232,"title":"Public Sport Channel","username":"sportchannel451","type":"channel"},"from":{"id":399749369,"is_bot":false,"first_name":"Artem","username":"vertex451","language_code":"en"},"date":1705664974,"old_chat_member":{"user":{"id":6845534569,"is_bot":true,"first_name":"Advertiser","username":"advertiser_451_bot"},"status":"left"},"new_chat_member":{"user":{"id":6845534569,"is_bot":true,"first_name":"Advertiser","username":"advertiser_451_bot"},"status":"administrator","can_be_edited":false,"can_manage_chat":true,"can_change_info":false,"can_post_messages":true,"can_edit_messages":false,"can_delete_messages":false,"can_invite_users":false,"can_restrict_members":true,"can_promote_members":false,"can_manage_video_chats":false,"can_post_stories":false,"can_edit_stories":false,"can_delete_stories":false,"is_anonymous":false,"can_manage_voice_chats":false}}}]}
// Remove bot from admins
// {"ok":true,"result":[{"update_id":632156491, "my_chat_member":{"chat":{"id":-1002025237232,"title":"Public Sport Channel","username":"sportchannel451","type":"channel"},"from":{"id":399749369,"is_bot":false,"first_name":"Artem","username":"vertex451","language_code":"en"},"date":1705664660,"old_chat_member":{"user":{"id":6845534569,"is_bot":true,"first_name":"Advertiser","username":"advertiser_451_bot"},"status":"administrator","can_be_edited":false,"can_manage_chat":true,"can_change_info":false,"can_post_messages":true,"can_edit_messages":false,"can_delete_messages":false,"can_invite_users":false,"can_restrict_members":true,"can_promote_members":false,"can_manage_video_chats":false,"can_post_stories":false,"can_edit_stories":false,"can_delete_stories":false,"is_anonymous":false,"can_manage_voice_chats":false},"new_chat_member":{"user":{"id":6845534569,"is_bot":true,"first_name":"Advertiser","username":"advertiser_451_bot"},"status":"left"}}}]}

func addNavigationButtons(msg tgbotapi.MessageConfig, buttons []tgbotapi.InlineKeyboardButton) tgbotapi.MessageConfig {
	if buttons == nil {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("« Back"), Back),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Start over ↺"), Start),
			),
		)
	} else {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			buttons,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("« Back"), Back),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Start over ↺"), Start),
			),
		)
	}

	return msg
}
