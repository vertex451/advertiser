package transport

import (
	"advertiser/shared/pkg/service/constants"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type CallBackQueryParams struct {
	ChatID   int64
	Page     string
	Variable string
}

func ParseCallBackQuery(query *tgbotapi.CallbackQuery) CallBackQueryParams {
	parsed := strings.Split(query.Data, "/")

	res := CallBackQueryParams{
		ChatID: query.Message.Chat.ID,
		Page:   parsed[0],
	}

	if len(parsed) > 1 {
		res.Variable = parsed[1]
	}

	return res
}

func AddNavigationButtons(msg tgbotapi.MessageConfig, buttons []tgbotapi.InlineKeyboardButton) tgbotapi.MessageConfig {
	if buttons == nil {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("« Back"), constants.Back),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Start over ↺"), constants.Start),
			),
		)
	} else {
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			buttons,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("« Back"), constants.Back),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Start over ↺"), constants.Start),
			),
		)
	}

	return msg
}
