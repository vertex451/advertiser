package transport

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/transport"
	"advertiser/shared/pkg/service/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func (s *Service) editTopicsPrompt(respondTo int64, rawChannelID string, storeInitialChannelData bool) types.CustomMessage {
	channelID, err := strconv.ParseInt(rawChannelID, 10, 64)
	if err != nil {
		zap.L().Panic("failed to parse channelID into int64 at editTopicsPrompt",
			zap.Error(err), zap.Any("rawChannelID", rawChannelID))
	}

	s.setState(respondTo, stateData{
		state:                   StateEditTopics,
		channelID:               channelID,
		storeInitialChannelData: storeInitialChannelData,
	})

	msgText := fmt.Sprintf(`
Оберіть топіки зі списку (напишіть їх через кому):
%s
`, strings.Join(s.uc.AllTopics(), ", "))

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, msgText),
		nil,
		false,
		false,
		false,
	)
}

func (s *Service) editChannelLocationPrompt(respondTo int64, rawChannelID string) types.CustomMessage {
	var locationButtons []tgbotapi.InlineKeyboardButton
	for locationCode, locationUkrName := range constants.Locations {
		data := fmt.Sprintf("%s/%s/%s", SetChannelLocation, rawChannelID, locationCode)
		locationButtons = append(locationButtons, tgbotapi.NewInlineKeyboardButtonData(locationUkrName, data))
	}

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, "Виберіть основну локацію каналу(чим точніше тим краще):"),
		transport.MakeTwoButtonsInARow(locationButtons),
		false,
		false,
		false,
	)
}

func (s *Service) editCostPerMilePrompt(respondTo int64, rawChannelID string) types.CustomMessage {
	channelID, err := strconv.ParseInt(rawChannelID, 10, 64)
	if err != nil {
		zap.L().Panic("failed to parse channelID into int64 at editTopicsPrompt",
			zap.Error(err), zap.Any("rawChannelID", rawChannelID))
	}

	s.setState(respondTo, stateData{
		state:     StateEditCostPerMile,
		channelID: channelID,
	})

	msgText := `Надішліть <b>бажану</b> ціну за тисячу переглядів в USD(Нариклад: 5.15). 
Канали з доступнішою ціною мають вищий пріорітет в черзі на рекламу.
`

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, msgText),
		nil,
		false,
		false,
		false,
	)
}

func (s *Service) reportBugPrompt(respondTo int64) types.CustomMessage {
	s.setState(respondTo, stateData{
		state: StateReportBug,
	})

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, "Опишіть детально проблему:"),
		nil,
		false,
		false,
		false,
	)
}

func (s *Service) requestFeaturePrompt(respondTo int64) types.CustomMessage {
	s.setState(respondTo, stateData{
		state: StateRequestFeature,
	})

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, "Опишіть бажаний фукнціонал:"),
		nil,
		false,
		false,
		false,
	)
}
