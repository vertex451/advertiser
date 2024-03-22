package transport

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/pkg/service/transport"
	"advertiser/shared/pkg/service/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func (s *Service) NavigateToPage(params transport.CallBackQueryParams) types.CustomMessage {
	switch params.Page {
	case constants.Start:
		return s.start(params.UserID)
	case constants.Back:
		return s.back(params.UserID)

	case constants.AllTopics:
		return s.allTopics(params.UserID)
	case Moderate:
		return s.moderate(params.UserID)
	case MyChannels:
		return s.listMyChannels(params.UserID)
	case ListChannelsTopics:
		return s.listChannelTopics(params.UserID, params.Variable)
	case EditChannelsTopics:
		return s.editTopicsPrompt(params.UserID, params.Variable)
	case ModerateDetails:
		return s.getAdvertisementDetails(params.UserID, params.Variable)
	case constants.ViewAdMessage:
		return s.viewAdMessage(params.UserID, params.Variable)
	case PostNow:
		return s.moderationDecision(params.UserID, PostNow, params.Variable)
	case RejectAd:
		return s.moderationDecision(params.UserID, RejectAd, params.Variable)

	default:
		return s.start(params.UserID)
	}
}

func (s *Service) start(respondTo int64) types.CustomMessage {
	s.resetState(respondTo)

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, "Choose action:"),
		[][]tgbotapi.InlineKeyboardButton{
			{tgbotapi.NewInlineKeyboardButtonData("My channels", fmt.Sprintf("%s", MyChannels))},
			{tgbotapi.NewInlineKeyboardButtonData("Moderation", fmt.Sprintf("%s", Moderate))},
			{tgbotapi.NewInlineKeyboardButtonData("All topics", fmt.Sprintf("%s", constants.AllTopics))},
		},
		true,
		false,
	)
}

func (s *Service) back(respondTo int64) types.CustomMessage {
	state := s.getState(respondTo)
	if len(state.crumbs) <= 1 {
		return s.start(respondTo)
	}

	params := state.crumbs[len(state.crumbs)-2]
	state.crumbs = state.crumbs[:len(state.crumbs)-1]
	state.state = StateStart

	s.setState(respondTo, state)

	return s.NavigateToPage(params)
}

func (s *Service) allTopics(respondTo int64) types.CustomMessage {
	msgText := fmt.Sprintf(`
Supported topics:
%s
`, strings.Join(s.uc.AllTopics(), ", "))

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, msgText),
		nil,
		true,
		false,
	)
}

func (s *Service) listMyChannels(respondTo int64) types.CustomMessage {
	var msg tgbotapi.MessageConfig
	myChannels, err := s.uc.ListMyChannels(respondTo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg = tgbotapi.NewMessage(respondTo, "You don's have Advertiser bot in your channels")
		} else {
			zap.L().Error("failed to list channels", zap.Error(err))
			msg = tgbotapi.NewMessage(respondTo, "failed to list channels")
		}

		return types.NewCustomMessageConfig(
			msg,
			nil,
			true,
			false,
		)

	}

	var channelButtons []tgbotapi.InlineKeyboardButton
	for channelID, channelName := range myChannels {
		data := fmt.Sprintf("%s/%s", ListChannelsTopics, strconv.FormatInt(channelID, 10))
		channelButtons = append(channelButtons, tgbotapi.NewInlineKeyboardButtonData(channelName, data))
	}

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, "Select a channel:"),
		transport.MakeTwoButtonsInARow(channelButtons),
		true,
		false,
	)
}

func (s *Service) listChannelTopics(respondTo int64, rawChannelID string) types.CustomMessage {
	channelID, err := strconv.ParseInt(rawChannelID, 10, 64)
	if err != nil {
		zap.L().Panic("failed to parse string to int64")
	}

	var msg tgbotapi.MessageConfig
	channelInfo, err := s.uc.GetChannelInfo(channelID)
	if err != nil {
		zap.L().Error("failed to list channel topics", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, "failed to list channel topics")
	} else {
		var topics []string
		for _, topic := range channelInfo.Topics {
			topics = append(topics, topic.ID)
		}
		text := fmt.Sprintf("%s topics: %s", channelInfo.Title, strings.Join(topics, ", "))
		msg = tgbotapi.NewMessage(respondTo, text)
	}

	buttons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("Edit %s topics", channelInfo.Title),
			fmt.Sprintf("%s/%s", EditChannelsTopics, strconv.FormatInt(channelID, 10))),
	)

	return types.NewCustomMessageConfig(
		msg,
		[][]tgbotapi.InlineKeyboardButton{buttons},
		true,
		false,
	)
}

func (s *Service) editChannelTopics(respondTo, channelID int64, topics []string) types.CustomMessage {
	s.resetState(channelID)

	var msg tgbotapi.MessageConfig

	var normalizedTopics []string
	for _, topic := range topics {
		normalizedTopics = append(normalizedTopics, strings.ToLower(strings.TrimSpace(topic)))
	}

	err := s.uc.UpdateChannelTopics(channelID, normalizedTopics)
	if err != nil {
		zap.L().Error("failed to update channel topics", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("failed to update channel topics. Error: %v", err))
	} else {
		text := fmt.Sprintf("Topics changed! New channel topics: %s", strings.Join(normalizedTopics, ", "))
		msg = tgbotapi.NewMessage(respondTo, text)
	}

	return types.NewCustomMessageConfig(
		msg,
		nil,
		true,
		false)

}

func (s *Service) moderate(id int64) types.CustomMessage {
	ads, err := s.uc.GetAdsToModerateByUserID(id)
	if err != nil {
		zap.L().Error("failed to get ads to moderate", zap.Error(err))
	}

	var msg tgbotapi.MessageConfig
	var rows [][]tgbotapi.InlineKeyboardButton
	if len(ads) == 0 {
		msg = tgbotapi.NewMessage(id, "No ads to moderate")
	} else {
		msg = tgbotapi.NewMessage(id, "Select an advertisement to moderate:")
		for _, entry := range ads {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s - %s (cpv: %v)", entry.Channel.Title, entry.Advertisement.Name, entry.Advertisement.CostPerView),
				fmt.Sprintf("%s/%s", ModerateDetails, entry.ID),
			)))
		}
	}

	return types.NewCustomMessageConfig(
		msg,
		rows,
		true,
		false,
	)

}

func (s *Service) getAdvertisementDetails(chatID int64, advertisementChannelID string) types.CustomMessage {
	advertisementChannel, err := s.uc.GetAdChanDetails(advertisementChannelID)
	if err != nil {
		zap.L().Error("failed to get ad details", zap.Error(err))
		return nil
	}

	var msg tgbotapi.MessageConfig
	msg = tgbotapi.NewMessage(chatID, fmt.Sprintf(
		`
Target channel: %s (@%s)
Advertisement details:
- Name: %s
- Cost per view: %v USD
`,
		advertisementChannel.Channel.Title,
		advertisementChannel.Channel.Handle,
		advertisementChannel.Advertisement.Name,
		advertisementChannel.Advertisement.CostPerView,
	))

	rows := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData(
			"View Advertisement message",
			fmt.Sprintf("%s/%s", constants.ViewAdMessage, advertisementChannel.ID),
		)},
	}

	return types.NewCustomMessageConfig(
		msg,
		rows,
		true,
		true,
	)
}

func (s *Service) viewAdMessage(chatID int64, adChanID string) types.CustomMessage {
	ad, err := s.uc.GetAdMessageByAdChanID(uuid.FromStringOrNil(adChanID))
	if err != nil {
		zap.L().Error("failed to get advertisement details", zap.Error(err))

		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to get advertisement details. Error: %v", err)),
			nil,
			true,
			false,
		)
	}

	msg := transport.ComposeAdMessage(
		chatID,
		*ad,
		[][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s", "Post Now"),
					fmt.Sprintf("%s/%s", PostNow, adChanID),
				),
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s", "Reject"),
					fmt.Sprintf("%s/%s", RejectAd, adChanID),
				),
			},
		},
		true,
		false,
	)

	return msg
}

func (s *Service) moderationDecision(respondTo int64, decision string, adChanID string) types.CustomMessage {
	var err error
	var msg tgbotapi.MessageConfig

	switch decision {
	case PostNow:
		err = s.PostAdvertisement(adChanID)
		if err != nil {
			msg = tgbotapi.NewMessage(respondTo, "Failed to post an advertisement")
		} else {
			msg = tgbotapi.NewMessage(respondTo, "Posted!")
		}
	case RejectAd:
		err = s.uc.UpdateAdChanEntry(models.AdvertisementChannel{
			ID:     uuid.FromStringOrNil(adChanID),
			Status: models.AdChanRejected,
		})
		if err != nil {
			msg = tgbotapi.NewMessage(respondTo, "Failed to reject an advertisement")
			zap.L().Error("failed to update advertisement status", zap.Error(err))
		} else {
			s.setState(respondTo, stateData{
				state:    StateWaitForRejectReason,
				adChanID: adChanID,
			})

			return types.NewCustomMessageConfig(
				tgbotapi.NewMessage(respondTo, "Please provide a reason for rejection:"),
				nil,
				false,
				false,
			)

		}
	}

	return types.NewCustomMessageConfig(
		msg,
		nil,
		true,
		false,
	)
}

func (s *Service) saveRejectionReason(respondTo int64, adChanID, reason string) types.CustomMessage {
	s.resetState(respondTo)
	err := s.uc.UpdateAdChanEntry(models.AdvertisementChannel{
		ID:              uuid.FromStringOrNil(adChanID),
		RejectionReason: reason,
	})
	if err != nil {
		zap.L().Error("failed to update advertisement status", zap.Error(err))
	}

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, "Thank you! We will review this."),
		nil,
		true,
		false,
	)
}
