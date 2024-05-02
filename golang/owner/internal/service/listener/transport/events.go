package transport

import (
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/pkg/service/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
)

func (s *Service) handleBotIsAddedToAdminsEvent(myChatMember *tgbotapi.ChatMemberUpdated) types.CustomMessage {
	errMsg := s.checkBotPermissions(myChatMember)
	if errMsg != nil {
		return errMsg
	}

	if s.chanelExist(myChatMember) {
		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(myChatMember.From.ID, `Потрібні дозволи отримані!
Не забудьте актуалізувати топіки, локацію та ціну за тисячу переглядів`),
			nil,
			false,
			false,
			true,
		)
	}

	errMsg, admins := s.getChatAdmins(myChatMember)
	if errMsg != nil {
		return errMsg
	}

	errMsg, membersCount := s.getChatMembersCount(myChatMember)
	if errMsg != nil {
		return errMsg
	}

	errMsg = s.storeInitialChatData(myChatMember, admins, membersCount)
	if errMsg != nil {
		return errMsg
	}

	return s.setInitialBotData(myChatMember)
}

func (s *Service) checkBotPermissions(myChatMember *tgbotapi.ChatMemberUpdated) (errMsg types.CustomMessage) {
	if !botHasNeededPermissions(myChatMember) {
		msg := tgbotapi.NewMessage(
			myChatMember.From.ID,
			fmt.Sprintf(`В каналі %s(@%s) недостатньо дозволів для роботи бота.
Бот потребує дозволів на публікацію та видалення повідомлень.
`,
				myChatMember.Chat.Title,
				myChatMember.Chat.UserName,
			),
		)

		return types.NewCustomMessageConfig(
			msg,
			nil,
			false,
			false,
			true,
		)

	}

	return nil
}

func botHasNeededPermissions(myChatMember *tgbotapi.ChatMemberUpdated) bool {
	return myChatMember.NewChatMember.CanPostMessages && myChatMember.NewChatMember.CanDeleteMessages
}

func (s *Service) chanelExist(myChatMember *tgbotapi.ChatMemberUpdated) bool {
	_, err := s.uc.GetChannelInfo(myChatMember.Chat.ID)
	if err != nil {
		return false
	}

	return true
}

func (s *Service) getChatAdmins(myChatMember *tgbotapi.ChatMemberUpdated) (types.CustomMessage, []tgbotapi.ChatMember) {
	admins, err := s.tgBotApi.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: myChatMember.Chat.ID,
		},
	})
	if err != nil {
		zap.L().Error("failed to get admin list, try to dismiss bot from admin and add it again.", zap.Error(err))
		msg := tgbotapi.NewMessage(
			myChatMember.From.ID,
			fmt.Sprintf("Не вдалося отримати список адміністраторів в %s, видаліть бота і додайте його ще раз.", myChatMember.Chat.Title),
		)

		return types.NewCustomMessageConfig(
			msg,
			nil,
			false,
			false,
			false,
		), nil
	}

	return nil, admins
}

func (s *Service) getChatMembersCount(myChatMember *tgbotapi.ChatMemberUpdated) (types.CustomMessage, int) {
	membersCount, err := s.tgBotApi.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: myChatMember.Chat.ID,
		},
	})
	if err != nil {
		zap.L().Error("failed to get members count", zap.Error(err))

		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("Failed to get members count")),
			nil,
			false,
			false,
			false,
		), 0
	}

	return nil, membersCount
}

func (s *Service) storeInitialChatData(myChatMember *tgbotapi.ChatMemberUpdated, admins []tgbotapi.ChatMember, membersCount int) types.CustomMessage {
	err := s.uc.StoreInitialChannelData(admins, models.Channel{
		ID:          myChatMember.Chat.ID,
		Description: myChatMember.Chat.Description,
		Handle:      myChatMember.Chat.UserName,
		IsChannel:   myChatMember.Chat.IsChannel(),
		Title:       myChatMember.Chat.Title,
		Subscribers: membersCount,
	})

	if err != nil {
		zap.L().Error("failed to store channel initial data", zap.Error(err))
		msg := tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("Не вдалося зберегти початкові дані каналу. Помилка: %v", err))

		return types.NewCustomMessageConfig(
			msg,
			nil,
			false,
			false,
			false,
		)
	}

	return nil
}

func (s *Service) setInitialBotData(myChatMember *tgbotapi.ChatMemberUpdated) types.CustomMessage {
	msgText := tgbotapi.NewMessage(
		myChatMember.From.ID,
		fmt.Sprintf(
			`Бот успішно доданий до %s! 
Для активації потрібні ще деякі дані.`,
			myChatMember.Chat.Title),
	)

	s.tgBotApi.Send(types.NewCustomMessageConfig(
		msgText,
		nil,
		false,
		false,
		false,
	))

	return s.editTopicsPrompt(myChatMember.From.ID, strconv.FormatInt(myChatMember.Chat.ID, 10), true)
}

func (s *Service) handleBotIsRemovedFromAdminsEvent(myChatMember *tgbotapi.ChatMemberUpdated) types.CustomMessage {
	var msg tgbotapi.MessageConfig
	err := s.uc.DeleteChannel(myChatMember.Chat.ID)
	if err != nil {
		zap.L().Error("failed to delete channel", zap.Error(err))
		msg = tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("failed to delete channel. Error: %v", err))
	} else {
		zap.L().Info("Advertiser bot is removed from channel", zap.String("chat", myChatMember.Chat.UserName))
		msg = tgbotapi.NewMessage(myChatMember.From.ID, fmt.Sprintf("Бот вилучений з %s", myChatMember.Chat.Title))
	}

	return types.NewCustomMessageConfig(
		msg,
		nil,
		false,
		false,
		true,
	)
}
