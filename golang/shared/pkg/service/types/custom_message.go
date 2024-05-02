package types

import (
	"advertiser/shared/pkg/service/constants"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CustomMessage interface {
	tgbotapi.Chattable
	SetReplyMarkup(markup interface{})
	SkipDeletionOfPrevMsg() bool
}

func NewCustomMessageConfig(
	msgConfig tgbotapi.MessageConfig,
	rows [][]tgbotapi.InlineKeyboardButton,
	addNavigation bool,
	reportBug bool,
	skipDeletion bool,
) CustomMessage {
	msg := &CustomMessageConfig{
		MessageConfig: msgConfig,
		skipDeletion:  skipDeletion,
	}

	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true

	if addNavigation {
		return AddNavigationButtons(msg, rows, reportBug)
	} else {
		if len(rows) > 0 {
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
		}
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
		return AddNavigationButtons(msg, rows, false)
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

func (msg *CustomMessageConfig) SetReplyMarkup(markup interface{}) {
	msg.ReplyMarkup = markup
}

func (msg *CustomMessageConfig) SkipDeletionOfPrevMsg() bool {
	return msg.skipDeletion
}

func (msg *CustomPhotoConfig) SetReplyMarkup(markup interface{}) {
	msg.ReplyMarkup = markup
}

func (msg *CustomPhotoConfig) SkipDeletionOfPrevMsg() bool {
	return msg.skipDeletion
}

func AddNavigationButtons(msg CustomMessage, rows [][]tgbotapi.InlineKeyboardButton, reportBug bool) CustomMessage {
	var replyMarkup tgbotapi.InlineKeyboardMarkup

	if rows == nil {
		replyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("« Назад"), constants.Back),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Заново ↺"), constants.Start),
			),
		)
	} else {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("« Назад"), constants.Back),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Заново ↺"), constants.Start),
		))
		replyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	}

	if reportBug {
		replyMarkup.InlineKeyboard = append(replyMarkup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Повідомити про помилку 🐞"), constants.ReportBug),
		))
	}

	msg.SetReplyMarkup(replyMarkup)

	return msg
}
