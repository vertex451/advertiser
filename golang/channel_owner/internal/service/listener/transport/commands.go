package transport

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/pkg/service/transport"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
)

func (s *Transport) handleCommand(update tgbotapi.Update) *tgbotapi.MessageConfig {
	switch update.Message.Command() {
	case constants.Start:
		return s.start(update.Message.Chat.ID)
	case constants.AllTopics:
		return s.allTopics(update.Message.Chat.ID)
	case Moderate:
		return s.moderate(update.Message.Chat.ID)
	}
	return nil
}

func (s *Transport) start(respondTo int64) *tgbotapi.MessageConfig {
	s.resetState(respondTo)

	msg := transport.AddNavigationButtons(
		tgbotapi.NewMessage(respondTo, "Choose action:"),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("List my channels with bot present", fmt.Sprintf("%s", MyChannels)),
			tgbotapi.NewInlineKeyboardButtonData("List all topics", fmt.Sprintf("%s", constants.AllTopics)),
		))

	return &msg
}

func (s *Transport) back(respondTo int64) *tgbotapi.MessageConfig {
	state := s.getState(respondTo)
	if len(state.crumbs) <= 1 {
		return s.start(respondTo)
	}

	params := state.crumbs[len(state.crumbs)-2]
	state.crumbs = state.crumbs[:len(state.crumbs)-1]

	s.setState(respondTo, state)

	return s.NavigateToPage(params)
}

func (s *Transport) allTopics(respondTo int64) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf(`
Supported topics:
%s
`, strings.Join(s.uc.AllTopics(), ", ")))
	msg = transport.AddNavigationButtons(msg, nil)

	return &msg
}

func (s *Transport) listMyChannels(respondTo int64) *tgbotapi.MessageConfig {
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

		return &msg
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

	return &msg
}

func (s *Transport) listChannelTopics(responseTo int64, rawChannelID string) *tgbotapi.MessageConfig {
	channelID, err := strconv.ParseInt(rawChannelID, 10, 64)
	if err != nil {
		zap.L().Panic("failed to parse string to int64")
	}

	var msg tgbotapi.MessageConfig
	channelInfo, err := s.uc.GetChannelInfo(channelID)
	if err != nil {
		zap.L().Error("failed to list channel topics", zap.Error(err))
		msg = tgbotapi.NewMessage(responseTo, "failed to list channel topics")
	} else {
		var topics []string
		for _, topic := range channelInfo.Topics {
			topics = append(topics, topic.ID)
		}
		text := fmt.Sprintf("%s topics: %s", channelInfo.Title, strings.Join(topics, ", "))
		msg = tgbotapi.NewMessage(responseTo, text)
	}

	buttons := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("Edit %s topics", channelInfo.Title),
			fmt.Sprintf("%s/%s", EditChannelsTopics, strconv.FormatInt(channelID, 10))),
	)
	msg = transport.AddNavigationButtons(msg, buttons)

	return &msg
}

func (s *Transport) editChannelTopics(responseTo, channelID int64, topics []string) *tgbotapi.MessageConfig {
	s.resetState(channelID)

	var msg tgbotapi.MessageConfig

	var normalizedTopics []string
	for _, topic := range topics {
		normalizedTopics = append(normalizedTopics, strings.ToLower(strings.TrimSpace(topic)))
	}

	err := s.uc.UpdateChannelTopics(channelID, normalizedTopics)
	if err != nil {
		zap.L().Error("failed to update channel topics", zap.Error(err))
		msg = tgbotapi.NewMessage(responseTo, fmt.Sprintf("failed to update channel topics. Error: %v", err))
	} else {
		text := fmt.Sprintf("Topics changed! New channel topics: %s", strings.Join(normalizedTopics, ", "))
		msg = tgbotapi.NewMessage(responseTo, text)
	}

	msg = transport.AddNavigationButtons(msg, nil)

	return &msg
}

func (s *Transport) moderate(id int64) *tgbotapi.MessageConfig {
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
			fmt.Println("### ads", entry)
			channelButtons = append(channelButtons,
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s (cpv: %v)", entry.ChannelTitle, entry.AdCostPerView),
					fmt.Sprintf("%s/%s", ModerateDetails, entry.ID),
				),
			)
		}
	}

	msg = transport.AddNavigationButtons(msg, channelButtons)

	return &msg
}

func (s *Transport) GetAdvertisementDetails(chatID int64, advertisementChannelID string) *tgbotapi.MessageConfig {
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
- Name: %s
- Cost per view: %v
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

	return &msg
}

func (s *Transport) moderationDecision(respondTo int64, decision string, adChanID string) *tgbotapi.MessageConfig {
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
		err = s.uc.UpdateAdChanStatus(adChanID, models.AdChanRejected)
		if err != nil {
			msg = tgbotapi.NewMessage(respondTo, "Failed to reject an advertisement")
			zap.L().Error("failed to update advertisement status", zap.Error(err))
		} else {
			msg = tgbotapi.NewMessage(respondTo, "Rejected!")
		}
	}

	msg = transport.AddNavigationButtons(msg, nil)

	return &msg
}
