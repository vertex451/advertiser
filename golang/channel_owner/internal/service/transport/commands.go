package transport

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
)

func (t *Transport) start(respondTo int64) *tgbotapi.MessageConfig {
	t.resetState(respondTo)

	msg := addNavigationButtons(
		tgbotapi.NewMessage(respondTo, "Choose action:"),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("List my channels with bot present", fmt.Sprintf("%s", MyChannels)),
			tgbotapi.NewInlineKeyboardButtonData("List all topics", fmt.Sprintf("%s", AllTopics)),
		))

	return &msg
}

func (t *Transport) back(respondTo int64) *tgbotapi.MessageConfig {
	state := t.getState(respondTo)
	if len(state.crumbs) <= 1 {
		return t.start(respondTo)
	}

	params := state.crumbs[len(state.crumbs)-2]
	state.crumbs = state.crumbs[:len(state.crumbs)-1]

	t.setState(respondTo, state)

	return t.NavigateToPage(params)
}

func (t *Transport) allTopics(respondTo int64) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf(`
Supported topics:
%s
`, strings.Join(t.uc.AllTopics(), ", ")))
	msg = addNavigationButtons(msg, nil)

	return &msg
}

func (t *Transport) listMyChannels(respondTo int64) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	myChannels, err := t.uc.ListMyChannels(respondTo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg = tgbotapi.NewMessage(respondTo, "You don't have Advertiser bot in your channels")
		} else {
			zap.L().Error("failed to list channels", zap.Error(err))
			msg = tgbotapi.NewMessage(respondTo, "failed to list channels")
		}

		msg = addNavigationButtons(msg, nil)

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
	msg = addNavigationButtons(msg, channelButtons)

	return &msg
}

func (t *Transport) listChannelTopics(responseTo int64, rawChannelID string) *tgbotapi.MessageConfig {
	channelID, err := strconv.ParseInt(rawChannelID, 10, 64)
	if err != nil {
		zap.L().Panic("failed to parse string to int64")
	}

	var msg tgbotapi.MessageConfig
	channelInfo, err := t.uc.GetChannelInfo(channelID)
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
	msg = addNavigationButtons(msg, buttons)

	return &msg
}

func (t *Transport) editChannelTopics(responseTo, channelID int64, topics []string) *tgbotapi.MessageConfig {
	t.resetState(channelID)

	var msg tgbotapi.MessageConfig

	var normalizedTopics []string
	for _, topic := range topics {
		normalizedTopics = append(normalizedTopics, strings.ToLower(strings.TrimSpace(topic)))
	}

	err := t.uc.UpdateChannelTopics(channelID, normalizedTopics)
	if err != nil {
		zap.L().Error("failed to update channel topics", zap.Error(err))
		msg = tgbotapi.NewMessage(responseTo, fmt.Sprintf("failed to update channel topics. Error: %v", err))
	} else {
		text := fmt.Sprintf("Topics changed! New channel topics: %s", strings.Join(normalizedTopics, ", "))
		msg = tgbotapi.NewMessage(responseTo, text)
	}

	msg = addNavigationButtons(msg, nil)

	return &msg
}
