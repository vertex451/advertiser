package transport

import (
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/pkg/service/transport"
	"fmt"
	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	StatusAdministrator = "administrator"
	StatusCreator       = "creator"
	StatusLeft          = "left"
)

func (s *Service) handleBotIsAddedToAdminsEvent(myChatMember *tgbotapi.ChatMemberUpdated) *transport.Msg {
	var err error
	var msg tgbotapi.MessageConfig
	if !botHasNeededPermissions(myChatMember) {
		msg = tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("Advertiser bot doesn's have needed permissions in channel %s.", myChatMember.Chat.Title))
		return &transport.Msg{
			Msg: msg,
		}
	}

	admins, err := s.tgBotApi.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
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

		return &transport.Msg{
			Msg: msg,
		}
	}

	membersCount, err := s.tgBotApi.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: myChatMember.Chat.ID,
		},
	})
	if err != nil {
		zap.L().Error("failed to get members count", zap.Error(err))
		msg = tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("Failed to get members count"))
		return &transport.Msg{
			Msg: msg,
		}
	}

	err = s.uc.StoreInitialChannelData(admins, models.Channel{
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

	return &transport.Msg{
		Msg: msg,
	}
}

func (s *Service) handleBotIsRemovedFromAdminsEvent(myChatMember *tgbotapi.ChatMemberUpdated) *transport.Msg {
	var msg tgbotapi.MessageConfig
	err := s.uc.DeleteChannel(myChatMember.Chat.ID)
	if err != nil {
		zap.L().Error("failed to delete channel", zap.Error(err))
		msg = tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("failed to delete channel. Error: %v", err))
	} else {
		zap.L().Info("Advertiser bot is removed from channel", zap.String("chat", myChatMember.Chat.UserName))
		msg = tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("Advertiser bot is removed from %s", myChatMember.Chat.Title))
	}

	return &transport.Msg{
		Msg: msg,
	}
}

func botHasNeededPermissions(myChatMember *tgbotapi.ChatMemberUpdated) bool {
	return myChatMember.NewChatMember.CanPostMessages && myChatMember.NewChatMember.CanDeleteMessages
}
