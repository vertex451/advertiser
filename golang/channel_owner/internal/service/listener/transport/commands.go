package transport

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/pkg/service/transport"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
)

func (s *Service) NavigateToPage(params transport.CallBackQueryParams) *transport.Msg {
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
		return s.GetAdvertisementDetails(params.UserID, params.Variable)
	case PostNow:
		return s.moderationDecision(params.UserID, PostNow, params.Variable)
	case RejectAd:
		return s.moderationDecision(params.UserID, RejectAd, params.Variable)

	default:
		return s.start(params.UserID)
	}
}

func (s *Service) start(respondTo int64) *transport.Msg {
	s.resetState(respondTo)

	var buttons []tgbotapi.InlineKeyboardButton
	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("My channels", fmt.Sprintf("%s", MyChannels)),
		tgbotapi.NewInlineKeyboardButtonData("Moderation", fmt.Sprintf("%s", Moderate)),
		tgbotapi.NewInlineKeyboardButtonData("All topics", fmt.Sprintf("%s", constants.AllTopics)),
	)

	msg := transport.AddNavigationButtons(
		tgbotapi.NewMessage(respondTo, "Choose action:"),
		buttons,
	)

	return &transport.Msg{
		Msg: msg,
	}
}

func (s *Service) back(respondTo int64) *transport.Msg {
	state := s.getState(respondTo)
	if len(state.crumbs) <= 1 {
		return s.start(respondTo)
	}

	params := state.crumbs[len(state.crumbs)-2]
	state.crumbs = state.crumbs[:len(state.crumbs)-1]

	s.setState(respondTo, state)

	return s.NavigateToPage(params)
}

func (s *Service) allTopics(respondTo int64) *transport.Msg {
	var msg tgbotapi.MessageConfig
	msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf(`
Supported topics:
%s
`, strings.Join(s.uc.AllTopics(), ", ")))
	msg = transport.AddNavigationButtons(msg, nil)

	return &transport.Msg{
		Msg: msg,
	}
}

func (s *Service) listMyChannels(respondTo int64) *transport.Msg {
	var msg tgbotapi.MessageConfig
	myChannels, err := s.uc.ListMyChannels(respondTo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg = tgbotapi.NewMessage(respondTo, "You don's have Advertiser bot in your channels")
		} else {
			zap.L().Error("failed to list channels", zap.Error(err))
			msg = tgbotapi.NewMessage(respondTo, "failed to list channels")
		}

		msg = transport.AddNavigationButtons(msg, nil)

		return &transport.Msg{
			Msg: msg,
		}
	}

	var channelButtons []tgbotapi.InlineKeyboardButton
	for channelID, channelName := range myChannels {
		data := fmt.Sprintf("%s/%s", ListChannelsTopics, strconv.FormatInt(channelID, 10))
		channelButtons = append(channelButtons, tgbotapi.NewInlineKeyboardButtonData(channelName, data))
	}

	sort.Slice(channelButtons, func(i, j int) bool {
		return channelButtons[i].Text < channelButtons[j].Text
	})

	msg = tgbotapi.NewMessage(respondTo, "Select a channel:")
	msg = transport.AddNavigationButtons(msg, channelButtons)

	return &transport.Msg{
		Msg: msg,
	}
}

func (s *Service) listChannelTopics(respondTo int64, rawChannelID string) *transport.Msg {
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
	msg = transport.AddNavigationButtons(msg, buttons)

	return &transport.Msg{
		Msg: msg,
	}
}

func (s *Service) editChannelTopics(respondTo, channelID int64, topics []string) *transport.Msg {
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

	msg = transport.AddNavigationButtons(msg, nil)

	return &transport.Msg{
		SkipDeletion: true,
		Msg:          msg,
	}
}

func (s *Service) moderate(id int64) *transport.Msg {
	ads, err := s.uc.GetAdsToModerateByUserID(id)
	if err != nil {
		zap.L().Error("failed to get ads to moderate", zap.Error(err))
	}

	var msg tgbotapi.MessageConfig
	var channelButtons []tgbotapi.InlineKeyboardButton
	if len(ads) == 0 {
		msg = tgbotapi.NewMessage(id, "No ads to moderate")
	} else {
		msg = tgbotapi.NewMessage(id, "Select an advertisement to moderate:")
		for _, entry := range ads {
			channelButtons = append(channelButtons,
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s (cpv: %v)", entry.ChannelTitle, entry.AdCostPerView),
					fmt.Sprintf("%s/%s", ModerateDetails, entry.ID),
				),
			)
		}
	}

	msg = transport.AddNavigationButtons(msg, channelButtons)

	return &transport.Msg{
		Msg: msg,
	}
}

func (s *Service) GetAdvertisementDetails(chatID int64, advertisementChannelID string) *transport.Msg {
	advertisementChannel, err := s.uc.GetAdChanDetails(advertisementChannelID)
	if err != nil {
		zap.L().Error("failed to get ad details", zap.Error(err))
		return nil
	}

	var msg tgbotapi.MessageConfig
	var channelButtons []tgbotapi.InlineKeyboardButton

	msg = tgbotapi.NewMessage(chatID, fmt.Sprintf(
		`
Target channel: %s
Advertisement details:
- ID: %s
- Cost per view: %v USD
- Message: %s
`,
		advertisementChannel.ChannelTitle,
		advertisementChannel.AdName,
		advertisementChannel.AdCostPerView,
		advertisementChannel.AdMessage,
	))
	channelButtons = append(channelButtons,
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s", "Post Now"),
			fmt.Sprintf("%s/%s", PostNow, advertisementChannel.ID),
		),
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s", "Reject"),
			fmt.Sprintf("%s/%s", RejectAd, advertisementChannel.ID),
		),
	)

	msg = transport.AddNavigationButtons(msg, channelButtons)

	return &transport.Msg{
		Msg: msg,
	}
}

func (s *Service) moderationDecision(respondTo int64, decision string, adChanID string) *transport.Msg {
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
			msg = tgbotapi.NewMessage(respondTo, "Please provide a reason for rejection")

			return &transport.Msg{
				Msg: msg,
			}
		}
	}

	msg = transport.AddNavigationButtons(msg, nil)

	return &transport.Msg{
		Msg: msg,
	}
}

func (s *Service) saveRejectionReason(respondTo int64, adChanID, reason string) *transport.Msg {
	s.resetState(respondTo)
	err := s.uc.UpdateAdChanEntry(models.AdvertisementChannel{
		ID:              uuid.FromStringOrNil(adChanID),
		RejectionReason: reason,
	})
	if err != nil {
		zap.L().Error("failed to update advertisement status", zap.Error(err))
	}

	msg := transport.AddNavigationButtons(
		tgbotapi.NewMessage(respondTo, "Thank you! We will review this."), nil,
	)

	return &transport.Msg{
		Msg: msg,
	}
}
