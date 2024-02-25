package transport

import (
	"advertiser/shared/pkg/service/constants"
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

func (s *Transport) moderationDecision(respondTo int64, decision string, adId string) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig

	switch decision {
	case ApproveAd:
		fmt.Println("### APPROVED", adId)
		msg = tgbotapi.NewMessage(respondTo, "Approved!")
	case RejectAd:
		fmt.Println("### REJECTED", adId)
		msg = tgbotapi.NewMessage(respondTo, "Rejected!")
	}

	return &msg
}