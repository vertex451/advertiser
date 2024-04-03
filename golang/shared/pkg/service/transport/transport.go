package transport

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/pkg/service/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"os"
	"sort"
	"strings"
)

var stringToHtmlMap = map[string]string{
	"bold":          "b",
	"code":          "code",
	"italic":        "i",
	"monospace":     "pre",
	"spoiler":       "tg-spoiler",
	"strikethrough": "strike",
	"underline":     "u",
}

type CallBackQueryParams struct {
	UserID   int64
	Page     string
	Variable string
}

func ParseCallBackQuery(query *tgbotapi.CallbackQuery) CallBackQueryParams {
	parsed := strings.Split(query.Data, "/")

	res := CallBackQueryParams{
		UserID: query.From.ID,
		Page:   parsed[0],
	}

	if len(parsed) > 1 {
		res.Variable = parsed[1]
	}

	return res
}

func MakeTwoButtonsInARow(buttons []tgbotapi.InlineKeyboardButton) [][]tgbotapi.InlineKeyboardButton {
	maxButtonsPerRow := 2

	sort.Slice(buttons, func(i, j int) bool {
		return buttons[i].Text < buttons[j].Text
	})

	var row []tgbotapi.InlineKeyboardButton
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, button := range buttons {
		row = append(row, button)

		// If the row is full, add it to the rows slice and start a new row
		if len(row) == maxButtonsPerRow {
			rows = append(rows, row)
			row = nil
		}
	}

	// Add the last row if it's not empty
	if len(row) > 0 {
		rows = append(rows, row)
	}

	return rows
}

func GetUserID(update tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.From.ID
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}

	if update.MyChatMember != nil {
		return update.MyChatMember.From.ID
	}

	return 0
}

// ComposeAdMessage applies the formatting on raw text.
// https://core.telegram.org/bots/api#formatting-options
func ComposeAdMessage(
	adChanChannelID int64,
	ad models.Advertisement,
	rows [][]tgbotapi.InlineKeyboardButton,
	addNavigation bool,
	skipDeletion bool,
) types.CustomMessage {
	msgText := composeText(ad)

	if ad.MsgImageURL != "" {
		return composeImageMessage(adChanChannelID, msgText, ad.MsgImageURL, rows, addNavigation, skipDeletion)
	} else {
		return composePlainMessage(adChanChannelID, msgText, rows, addNavigation, skipDeletion)
	}
}

func composeText(ad models.Advertisement) string {
	var composedMsg strings.Builder
	currentPos := 0

	// Sort entities by offset
	sortEntitiesByOffset(ad.MsgEntities)

	// Iterate over entities
	for _, entity := range ad.MsgEntities {
		// Append the text before the current entity
		composedMsg.WriteString(ad.MsgText[currentPos:entity.Offset])

		// Append the formatted text according to the entity type
		switch entity.Type {
		case "bold", "code", "italic", "pre", "spoiler", "strikethrough", "underline":
			composedMsg.WriteString(fmt.Sprintf("<%s>%s</%s>",
				stringToHtmlMap[entity.Type],
				ad.MsgText[entity.Offset:entity.Offset+entity.Length],
				stringToHtmlMap[entity.Type]),
			)
		case "text_link":
			composedMsg.WriteString(fmt.Sprintf("<a href=\"%s\">%s</a>", entity.URL, ad.MsgText[entity.Offset:entity.Offset+entity.Length]))
		default:
			composedMsg.WriteString(ad.MsgText[entity.Offset : entity.Offset+entity.Length])
		}

		currentPos = entity.Offset + entity.Length
	}

	// Append the remaining text after the last entity
	composedMsg.WriteString(ad.MsgText[currentPos:])

	return composedMsg.String()
}

func composeImageMessage(
	adChanChannelID int64,
	msgText, msgImageURL string,
	rows [][]tgbotapi.InlineKeyboardButton,
	addNavigation, skipDeletion bool,
) types.CustomMessage {
	imgBytes, err := os.ReadFile(msgImageURL)
	if err != nil {
		zap.L().Error("failed to load image from disk")
		return composePlainMessage(adChanChannelID, msgText, rows, addNavigation, skipDeletion)
	}

	msg := tgbotapi.NewPhoto(
		adChanChannelID,
		tgbotapi.FileBytes{
			Name:  msgImageURL,
			Bytes: imgBytes,
		})

	msg.Caption = msgText
	msg.ParseMode = tgbotapi.ModeHTML

	return types.NewCustomPhotoConfig(
		msg,
		rows,
		addNavigation,
		skipDeletion,
	)
}

func composePlainMessage(
	adChanChannelID int64,
	msgText string,
	rows [][]tgbotapi.InlineKeyboardButton,
	addNavigation, skipDeletion bool,
) types.CustomMessage {
	msg := tgbotapi.NewMessage(
		adChanChannelID,
		msgText,
	)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true

	return types.NewCustomMessageConfig(
		msg,
		rows,
		addNavigation,
		skipDeletion,
	)
}

func sortEntitiesByOffset(entities []models.MsgEntity) {
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Offset < entities[j].Offset
	})
}

func AddFooter(msg string) string {
	return fmt.Sprintf(
		`%s
	
- - - - -
Advertisement network "To Infinity and Beyond" @%s`,
		msg,
		constants.ChannelMonetizerBotName)

}
