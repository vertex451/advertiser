package transport

import (
	"advertiser/owner/internal/service/listener"
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/transport"
	"advertiser/shared/pkg/service/types"
	"advertiser/shared/tg_api"
	"advertiser/shared/tg_bot_api"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"strings"
	"sync"
)

// Commands:
const (
	EditChannelsTopics = "edit_channel_topics"
	ListChannelsTopics = "list_channels_topics"
	MyChannels         = "my_channels"

	Moderate        = "moderate"
	ModerateDetails = "moderate_details"
	PostNow         = "post_now"
	PostLater       = "post_later"
	RejectAd        = "reject_ad"
)

type BotState int

const (
	StateStart BotState = iota
	StateEditTopics
	StateWaitForRejectReason
)

type Service struct {
	uc       listener.UseCase
	tgBotApi tg_bot_api.TgBotApiProvider
	state    sync.Map // map[UserID]stateData
	cron     *cron.Cron
	tgApi    *tg_api.Service
	env      string
}

func New(uc listener.UseCase, tgBotApi tg_bot_api.TgBotApiProvider, env string) *Service {
	c := cron.New(
		cron.WithParser(
			cron.NewParser(
				cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
			)))
	c.Start()

	return &Service{
		tgBotApi: tgBotApi,
		uc:       uc,
		cron:     c,
		tgApi:    tg_api.New(),
		env:      env,
	}
}

func (s *Service) MonitorChannels() {
	var err error
	var sentMsg tgbotapi.Message
	var state stateData
	var userID int64
	updates := s.tgBotApi.GetUpdatesChan()
	for update := range updates {
		fmt.Printf("\n### update %+v\n\n", update)
		if update.CallbackQuery != nil {
			fmt.Printf("\n### update.CallbackQuery.Message %+v\n\n", update.CallbackQuery.Message)
		}
		if update.Message != nil {
			fmt.Printf("\n### update.Message %+v\n\n", update.Message)
		}

		responseMessage := s.handleUpdate(update)
		if responseMessage == nil {
			continue
		}

		userID = transport.GetUserID(update)

		state = s.getState(userID)
		if !responseMessage.SkipDeletion() && state.lastMsgID != 0 {
			deleteMsg := tgbotapi.NewDeleteMessage(userID, state.lastMsgID)
			_, err = s.tgBotApi.Send(deleteMsg)
			if err != nil {
				zap.L().Error("failed to delete message", zap.Error(err))
			}
		}

		//_, err = s.tgBotApi.Send(tgbotapi.NewMessage(userID, "⬇"))
		//if err != nil {
		//	zap.L().Error("failed to add empty line", zap.Error(err))
		//}

		sentMsg, err = s.tgBotApi.Send(responseMessage)
		if err != nil {
			zap.L().Error("failed to send message", zap.Error(err))
			continue
		}
		state.lastMsgID = sentMsg.MessageID
		s.state.Store(userID, state)
	}
}

func (s *Service) handleUpdate(update tgbotapi.Update) types.CustomMessage {
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

func (s *Service) handleCommand(update tgbotapi.Update) types.CustomMessage {
	userID := transport.GetUserID(update)
	switch update.Message.Command() {
	case constants.Start:
		return s.start(userID)
	case constants.AllTopics:
		return s.allTopics(userID)
	case Moderate:
		return s.moderate(userID)
	}
	return nil
}

func (s *Service) handleStateQuery(update tgbotapi.Update) types.CustomMessage {
	userID := transport.GetUserID(update)
	state := s.getState(userID)

	switch state.state {
	case StateEditTopics:
		topics := strings.Split(update.Message.Text, ",")
		return s.editChannelTopics(userID, state.channelID, topics)
	case StateWaitForRejectReason:
		return s.saveRejectionReason(userID, state.adChanID, update.Message.Text)
	default:
		return s.start(userID)
	}
}

func (s *Service) handleCallbackQuery(query *tgbotapi.CallbackQuery) types.CustomMessage {
	params := transport.ParseCallBackQuery(query)

	if params.Page != constants.Back {
		s.addCrumbs(params)
	}

	return s.NavigateToPage(params)
}

func (s *Service) handleUpdateEvent(update tgbotapi.Update) types.CustomMessage {
	switch update.MyChatMember.NewChatMember.Status {
	case constants.StatusAdministrator:
		return s.handleBotIsAddedToAdminsEvent(update.MyChatMember)
	case constants.StatusLeft:
		return s.handleBotIsRemovedFromAdminsEvent(update.MyChatMember)
	}

	return nil
}
