package transport

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/transport"
	"advertiser/shared/pkg/service/types"
	"advertiser/shared/pkg/storage"
	"advertiser/shared/tg_bot_api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"sync"
	"tg-bot/internal/service/bot_api"
)

const (
	TgBotDirectChatID int64 = 6406834985
)

// Commands:
const (
	AllTopicsWithCoverage = "all_topics_with_coverage"
	CreateCampaign        = "create_campaign"
	MyCampaigns           = "my_campaigns"
	CampaignDetails       = "campaign_details"

	CreateAd  = "create_ad"
	AdDetails = "ad_details"

	EditAd   = "edit_ad"
	DeleteAd = "delete_ad"
	RunAd    = "run_ad"
	PauseAd  = "pause_ad"
	FinishAd = "finish_ad"
)

type BotState int

const (
	StateStart BotState = iota
	StateSetCampaignName
	StateCreateAd
	StateCreateAdMessage
	//StateUpdateAd
)

type Transport struct {
	tgBotApi     tg_bot_api.TgBotApiProvider
	uc           bot_api.UseCase
	storage      storage.Provider
	updateConfig tgbotapi.UpdateConfig
	state        sync.Map // map[UserID]stateData
	env          string
}

func New(
	uc bot_api.UseCase,
	tgBotApi tg_bot_api.TgBotApiProvider,
	storage storage.Provider,
	env string,
) *Transport {
	zap.L().Info("Started Transport")

	return &Transport{
		tgBotApi: tgBotApi,
		uc:       uc,
		storage:  storage,
		env:      env,
	}
}

func (t *Transport) MonitorChannels() {
	var err error
	var sentMsg tgbotapi.Message
	var state stateData
	var userID int64
	updates := t.tgBotApi.GetUpdatesChan()
	for update := range updates {
		responseMessage := t.handleUpdate(update)
		if responseMessage == nil {
			continue
		}

		userID = transport.GetUserID(update)
		state = t.getState(userID)

		if !responseMessage.SkipDeletion() && state.lastMsgID != 0 {
			deleteMsg := tgbotapi.NewDeleteMessage(TgBotDirectChatID, state.lastMsgID)
			t.tgBotApi.Send(deleteMsg)
		}

		sentMsg, err = t.tgBotApi.Send(responseMessage)
		if err != nil {
			zap.L().Error("failed to send message", zap.Error(err))
			continue
		}
		state.lastMsgID = sentMsg.MessageID
		t.state.Store(userID, state)
	}
}

// TODO add rate limiter
func (t *Transport) handleUpdate(update tgbotapi.Update) types.CustomMessage {
	userID := transport.GetUserID(update)
	if update.Message != nil {
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case constants.Start:
				return t.start(userID)
			case constants.AllTopics:
				return t.allTopics(userID)
			}

		}

		if update.Message.Text != "" || update.Message.Caption != "" {
			return t.handleStateQuery(update)
		}
	}

	if update.CallbackQuery != nil {
		return t.handleCallbackQuery(update.CallbackQuery)
	}

	return nil
}
