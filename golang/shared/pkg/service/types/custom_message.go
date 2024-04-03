package types

import (
	"advertiser/shared/pkg/service/constants"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CustomMessage interface {
	tgbotapi.Chattable
	// GetMessageConfig() tgbotapi.MessageConfig
	//SetChatID(int64)
	SetReplyMarkup(markup interface{})
	SkipDeletion() bool
}

func NewCustomMessageConfig(
	msgConfig tgbotapi.MessageConfig,
	rows [][]tgbotapi.InlineKeyboardButton,
	addNavigation bool,
	skipDeletion bool,
) CustomMessage {
	msg := &CustomMessageConfig{
		MessageConfig: msgConfig,
		skipDeletion:  skipDeletion,
	}

	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true

	if addNavigation {
		return AddNavigationButtons(msg, rows)
	} else {
		return msg
	}
}

func NewCustomPhotoConfig(
	photoConfig tgbotapi.PhotoConfig,
	rows [][]tgbotapi.InlineKeyboardButton,
	addNavigation bool,
	skipDeletion bool,
) CustomMessage {
	msg := &CustomPhotoConfig{
		PhotoConfig:  photoConfig,
		skipDeletion: skipDeletion,
	}

	if addNavigation {
		return AddNavigationButtons(msg, rows)
	} else {
		return msg
	}
}

type CustomMessageConfig struct {
	tgbotapi.MessageConfig
	skipDeletion bool
}

type CustomPhotoConfig struct {
	tgbotapi.PhotoConfig
	skipDeletion bool
}

//func (msg *CustomMessageConfig) SetChatID(chatID int64) {
//	msg.ChatID = chatID
//}

func (msg *CustomMessageConfig) SetReplyMarkup(markup interface{}) {
	msg.ReplyMarkup = markup
}

func (msg *CustomMessageConfig) SkipDeletion() bool {
	return msg.skipDeletion
}

//func (msg *CustomMessageConfig) GetMessageConfig() tgbotapi.MessageConfig {
//	return msg.MessageConfig
//}

//func (msg *CustomPhotoConfig) SetChatID(chatID int64) {
//	msg.ChatID = chatID
//}

func (msg *CustomPhotoConfig) SetReplyMarkup(markup interface{}) {
	msg.ReplyMarkup = markup
}

func (msg *CustomPhotoConfig) SkipDeletion() bool {
	return msg.skipDeletion
}

//func (msg *CustomPhotoConfig) GetMessageConfig() tgbotapi.MessageConfig {
//	return tgbotapi.NewMessage(msg.ChatID, "Photo")
//}

func AddNavigationButtons(msg CustomMessage, rows [][]tgbotapi.InlineKeyboardButton) CustomMessage {
	if rows == nil {
		msg.SetReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("« Back"), constants.Back),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Start over ↺"), constants.Start),
			),
		))
	} else {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("« Back"), constants.Back),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Start over ↺"), constants.Start),
		))
		msg.SetReplyMarkup(tgbotapi.NewInlineKeyboardMarkup(rows...))
	}

	return msg
}
