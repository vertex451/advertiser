package transport

import (
	"advertiser/shared/pkg/repo/models"
	"fmt"
	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	StatusAdministrator = "administrator"
	StatusLeft          = "left"
)

func (t *Transport) handleUpdateEvent(update tgbotapi.Update) *tgbotapi.MessageConfig {
	switch update.MyChatMember.NewChatMember.Status {
	case StatusAdministrator:
		return t.botIsAddedToAdmins(update.MyChatMember)
	case StatusLeft:
		return t.botIsRemovedFromAdmins(update.MyChatMember)
	}

	return nil
}

func (t *Transport) botIsAddedToAdmins(myChatMember *tgbotapi.ChatMemberUpdated) *tgbotapi.MessageConfig {
	var err error
	var msg tgbotapi.MessageConfig
	if !myChatMember.NewChatMember.CanPostMessages {
		msg = tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("Advertiser bot doesn't have needed permissions in channel %s.", myChatMember.Chat.Title))
		return &msg
	}

	admins, err := t.tgBotApi.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: myChatMember.Chat.ID,
		},
	})
	if err != nil {
		zap.L().Error("failed to get admin list", zap.Error(err))
		msg = tgbotapi.NewMessage(
			myChatMember.From.ID,
			fmt.Sprintf("Failed to get admin list for %s, please check bots permissions", myChatMember.Chat.Title),
		)

		return &msg
	}

	membersCount, err := t.tgBotApi.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: myChatMember.Chat.ID,
		},
	})
	if err != nil {
		zap.L().Error("failed to get members count", zap.Error(err))
		msg = tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("Failed to get members count"))
		return &msg
	}

	err = t.uc.StoreInitialChannelData(admins, models.Channel{
		ID:          myChatMember.Chat.ID,
		Description: myChatMember.Chat.Description,
		Handle:      myChatMember.Chat.UserName,
		IsChannel:   myChatMember.Chat.IsChannel(),
		Title:       myChatMember.Chat.Title,
		Subscribers: membersCount,
	})

	if err != nil {
		zap.L().Error("failed to store channel initial data", zap.Error(err))
	}

	msg = tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("Advertiser bot was successfully added to %s", myChatMember.Chat.Title))

	return &msg
}

func (t *Transport) botIsRemovedFromAdmins(myChatMember *tgbotapi.ChatMemberUpdated) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	err := t.uc.DeleteChannel(myChatMember.Chat.ID)
	if err != nil {
		zap.L().Error("failed to delete channel", zap.Error(err))
		msg = tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("failed to delete channel. Error: %v", err))
	} else {
		zap.L().Info("Advertiser bot is removed from channel", zap.String("chat", myChatMember.Chat.UserName))
		msg = tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("Advertiser bot is removed from %s", myChatMember.Chat.Title))
	}

	return &msg
}
