package transport

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/pkg/service/transport"
	"advertiser/shared/pkg/service/types"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func (s *Service) NavigateToPage(params transport.CallBackQueryParams) types.CustomMessage {
	switch params.Page {
	case constants.Start:
		return s.start(params.UserID)
	case constants.Back:
		return s.back(params.UserID)

	case constants.Help:
		return s.help(params.UserID)
	case constants.ReportBug:
		return s.reportBugPrompt(params.UserID)
	case constants.RequestFeature:
		return s.requestFeaturePrompt(params.UserID)
	case constants.Support:
		return s.support(params.UserID)

	case constants.AllTopics:
		return s.allTopics(params.UserID)
	case Moderate:
		return s.moderate(params.UserID)
	case MyChannels:
		return s.listMyChannels(params.UserID)
	case ListChannelInfo:
		return s.listChannelTopics(params.UserID, params.Variable)
	case EditChannelTopics:
		return s.editTopicsPrompt(params.UserID, params.Variable, false)
	case EditChannelLocation:
		return s.editChannelLocationPrompt(params.UserID, params.Variable)
	case SetChannelLocation:
		return s.setChannelLocation(params.UserID, params.Variable, params.SecondVariable)
	case ModerateDetails:
		return s.getAdvertisementDetails(params.UserID, params.Variable)
	case constants.ViewAdMessage:
		return s.viewAdMessage(params.UserID, params.Variable)
	case PostNow:
		return s.moderationDecision(params.UserID, PostNow, params.Variable)
	case RejectAd:
		return s.moderationDecision(params.UserID, RejectAd, params.Variable)

	default:
		return s.start(params.UserID)
	}
}

func (s *Service) start(respondTo int64) types.CustomMessage {
	s.resetState(respondTo)

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, "Що шукаєте?"),
		[][]tgbotapi.InlineKeyboardButton{
			{tgbotapi.NewInlineKeyboardButtonData("Мої канали з ботом", fmt.Sprintf("%s", MyChannels))},
			{tgbotapi.NewInlineKeyboardButtonData("Пропозиції по рекламі", fmt.Sprintf("%s", Moderate))},
			{tgbotapi.NewInlineKeyboardButtonData("Всі топіки", fmt.Sprintf("%s", constants.AllTopics))},
			{tgbotapi.NewInlineKeyboardButtonData("Допомога", fmt.Sprintf("%s", constants.Help))},
		},
		true,
		false,
		false,
	)
}

func (s *Service) help(respondTo int64) types.CustomMessage {
	s.resetState(respondTo)

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, "Що там у вас?"),
		[][]tgbotapi.InlineKeyboardButton{
			{tgbotapi.NewInlineKeyboardButtonData("Повідомити про помилку 🐞", fmt.Sprintf("%s", constants.ReportBug))},
			{tgbotapi.NewInlineKeyboardButtonData("Запросити додатковий функціонал", fmt.Sprintf("%s", constants.RequestFeature))},
			{tgbotapi.NewInlineKeyboardButtonData("Звернутися в технічну підтримку", fmt.Sprintf("%s", constants.Support))},
		},

		true,
		false,
		false,
	)
}

func (s *Service) reportBug(respondTo int64, bugDescription string) types.CustomMessage {
	s.resetState(respondTo)

	var msg tgbotapi.MessageConfig

	err := s.uc.ReportBug(respondTo, bugDescription)
	if err != nil {
		zap.L().Error("failed to report bug", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, "Не вдалося відправити ваш звіт")
	} else {
		msg = tgbotapi.NewMessage(respondTo, "Дякую за ваш звіт!")
	}

	return types.NewCustomMessageConfig(
		msg,
		nil,
		true,
		false,
		false,
	)
}

func (s *Service) requestFeature(respondTo int64, featureDescription string) types.CustomMessage {
	s.resetState(respondTo)

	var msg tgbotapi.MessageConfig

	err := s.uc.RequestFeature(respondTo, featureDescription)
	if err != nil {
		zap.L().Error("failed to request feature", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, "Не вдалося відправити ваш запит")
	} else {
		msg = tgbotapi.NewMessage(respondTo, "Дякую за ваш запит!")
	}

	return types.NewCustomMessageConfig(
		msg,
		nil,
		true,
		false,
		false,
	)
}

func (s *Service) support(respondTo int64) types.CustomMessage {
	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, fmt.Sprintf("Надішліть повідомлення на адресу %s", constants.TolokaDigitalEmail)),
		nil,
		true,
		false,
		false,
	)
}

func (s *Service) back(respondTo int64) types.CustomMessage {
	state := s.getState(respondTo)
	if len(state.crumbs) <= 1 {
		return s.start(respondTo)
	}

	params := state.crumbs[len(state.crumbs)-2]
	state.crumbs = state.crumbs[:len(state.crumbs)-1]
	state.state = StateStart

	s.setState(respondTo, state)

	return s.NavigateToPage(params)
}

func (s *Service) allTopics(respondTo int64) types.CustomMessage {
	msgText := fmt.Sprintf(`
Наявні топіки:
%s
`, strings.Join(s.uc.AllTopics(), ", "))

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, msgText),
		nil,
		true,
		false,
		false,
	)
}

func (s *Service) listMyChannels(respondTo int64) types.CustomMessage {
	var msg tgbotapi.MessageConfig
	myChannels, err := s.uc.ListMyChannels(respondTo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg = tgbotapi.NewMessage(respondTo,
				`У вас ще не має бота в каналах.
Щоб користуватися ботом, додайте його як адміністратора до каналу з наступними дозволами:
1. Дозвіл на керування повідомленнями 3/3 
(не хвилюйтеся, бот буде публікувати рекламні повідомлення тільки після узгодження з вами)
`)
		} else {
			zap.L().Error("failed to list channels", zap.Error(err))
			msg = tgbotapi.NewMessage(respondTo, "Не вдалося отримати список каналів")
		}

		return types.NewCustomMessageConfig(
			msg,
			nil,
			true,
			false,
			false,
		)

	}

	var channelButtons []tgbotapi.InlineKeyboardButton
	for channelID, channelName := range myChannels {
		data := fmt.Sprintf("%s/%s", ListChannelInfo, strconv.FormatInt(channelID, 10))
		channelButtons = append(channelButtons, tgbotapi.NewInlineKeyboardButtonData(channelName, data))
	}

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, "Оберіть канал:"),
		transport.MakeTwoButtonsInARow(channelButtons),
		true,
		false,
		false,
	)
}

func (s *Service) listChannelTopics(respondTo int64, rawChannelID string) types.CustomMessage {
	channelID, err := strconv.ParseInt(rawChannelID, 10, 64)
	if err != nil {
		zap.L().Panic("failed to parse string to int64")
	}

	var msg tgbotapi.MessageConfig
	channelInfo, err := s.uc.GetChannelInfo(channelID)
	if err != nil {
		zap.L().Error("failed to get channel info", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, "Не вдалося отримати інформацію про канал")
	} else {
		var topics []string
		for _, topic := range channelInfo.Topics {
			topics = append(topics, topic.ID)
		}

		var text string
		if len(topics) == 0 {
			text = fmt.Sprintf("%s не має топіків, задайте їх щоб отримувати пропозиції по рекламі", channelInfo.Title)
		} else {
			text = fmt.Sprintf(`<b>%s</b>
Локація: %s
Топіки: %s`,
				channelInfo.Title,
				constants.Locations[channelInfo.Location],
				strings.Join(topics, ", "),
			)
		}
		msg = tgbotapi.NewMessage(respondTo, text)
	}

	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData(
				"Змінити локацію",
				fmt.Sprintf("%s/%s", EditChannelLocation, strconv.FormatInt(channelID, 10)),
			),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData(
				"Змінити топіки",
				fmt.Sprintf("%s/%s", EditChannelTopics, strconv.FormatInt(channelID, 10)),
			),
		},
	}

	return types.NewCustomMessageConfig(
		msg,
		buttons,
		true,
		false,
		false,
	)
}

func (s *Service) editChannelTopics(respondTo, channelID int64, topics []string) types.CustomMessage {
	state := s.getState(respondTo)

	s.resetState(respondTo)

	var msg tgbotapi.MessageConfig

	var normalizedTopics []string
	for _, topic := range topics {
		normalizedTopics = append(normalizedTopics, strings.ToLower(strings.TrimSpace(topic)))
	}

	err := s.uc.UpdateChannelTopics(channelID, normalizedTopics)
	if err != nil {
		zap.L().Error("failed to update channel topics", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Не вдалося оновити топіки каналу. Помилка: %v", err))
	} else {
		text := fmt.Sprintf("Оновлено, поточні топіки: %s", strings.Join(normalizedTopics, ", "))
		msg = tgbotapi.NewMessage(respondTo, text)
	}

	if state.storeInitialChannelData {
		s.tgBotApi.Send(types.NewCustomMessageConfig(
			msg,
			nil,
			false,
			false,
			false,
		))

		return s.editChannelLocationPrompt(respondTo, strconv.FormatInt(channelID, 10))
	} else {
		return types.NewCustomMessageConfig(
			msg,
			nil,
			true,
			false,
			false,
		)
	}
}

func (s *Service) setChannelLocation(respondTo int64, rawChannelID, rawLocation string) types.CustomMessage {
	s.resetState(respondTo)

	channelID, err := strconv.ParseInt(rawChannelID, 10, 64)
	if err != nil {
		zap.L().Error("failed to parse string to int64 in setChannelLocation")
	}

	location, err := strconv.Atoi(rawLocation)
	if err != nil {
		zap.L().Error("failed to parse string to int in setChannelLocation")
	}

	var msg tgbotapi.MessageConfig

	err = s.uc.UpdateChannelLocation(channelID, constants.Location(location))
	if err != nil {
		zap.L().Error("failed to update channel location", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Не вдалося оновити локацію каналу. Помилка: %v", err))
	} else {
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Оновлено, поточна локація - %s", constants.Locations[constants.Location(location)]))
	}

	return types.NewCustomMessageConfig(
		msg,
		nil,
		true,
		false,
		false,
	)
}

func (s *Service) moderate(id int64) types.CustomMessage {
	ads, err := s.uc.GetAdsToModerateByUserID(id)
	if err != nil {
		zap.L().Error("failed to get ads to moderate", zap.Error(err))
	}

	var msg tgbotapi.MessageConfig
	var rows [][]tgbotapi.InlineKeyboardButton
	if len(ads) == 0 {
		msg = tgbotapi.NewMessage(id, "Пропозиції по рекламі ще не настоялися, перевірте пізніше")
	} else {
		msg = tgbotapi.NewMessage(id, "Виберіть рекламу і винесіть рішення:")
		for _, entry := range ads {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%s - %s (cpv: %v)", entry.Channel.Title, entry.Advertisement.Name, entry.Advertisement.CostPerMile),
				fmt.Sprintf("%s/%s", ModerateDetails, entry.ID),
			)))
		}
	}

	return types.NewCustomMessageConfig(
		msg,
		rows,
		true,
		false,
		false,
	)

}

func (s *Service) getAdvertisementDetails(chatID int64, advertisementChannelID string) types.CustomMessage {
	advertisementChannel, err := s.uc.GetAdChanDetails(advertisementChannelID)
	if err != nil {
		zap.L().Error("failed to get ad details", zap.Error(err))
		return nil
	}

	var msg tgbotapi.MessageConfig
	msg = tgbotapi.NewMessage(chatID, fmt.Sprintf(
		`
Цільовий канал: %s (@%s)
Деталі реклами:
- Назва: %s
- Cost per mile: %v USD
`,
		advertisementChannel.Channel.Title,
		advertisementChannel.Channel.Handle,
		advertisementChannel.Advertisement.Name,
		advertisementChannel.Advertisement.CostPerMile,
	))

	rows := [][]tgbotapi.InlineKeyboardButton{
		{tgbotapi.NewInlineKeyboardButtonData(
			"Дивитись рекламне повідомлення",
			fmt.Sprintf("%s/%s", constants.ViewAdMessage, advertisementChannel.ID),
		)},
	}

	return types.NewCustomMessageConfig(
		msg,
		rows,
		true,
		false,
		true,
	)
}

func (s *Service) viewAdMessage(chatID int64, adChanID string) types.CustomMessage {
	ad, err := s.uc.GetAdMessageByAdChanID(uuid.FromStringOrNil(adChanID))
	if err != nil {
		zap.L().Error("failed to get advertisement details", zap.Error(err))

		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(chatID, fmt.Sprintf("Не вдалося показати рекламне повідомлення. Помикла: %v", err)),
			nil,
			true,
			false,
			false,
		)
	}

	msg := transport.ComposeAdMessage(
		chatID,
		*ad,
		[][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s", "Опублікувати зараз"),
					fmt.Sprintf("%s/%s", PostNow, adChanID),
				),
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s", "Відхилити"),
					fmt.Sprintf("%s/%s", RejectAd, adChanID),
				),
			},
		},
		true,
		false,
	)

	return msg
}

func (s *Service) moderationDecision(respondTo int64, decision string, adChanID string) types.CustomMessage {
	var err error
	var msg tgbotapi.MessageConfig

	switch decision {
	case PostNow:
		err = s.PostAdvertisement(adChanID)
		if err != nil {
			msg = tgbotapi.NewMessage(respondTo, "Не вдалося опублікувати рекламу")
		} else {
			msg = tgbotapi.NewMessage(respondTo, "Опубліковано! Готуйте мішки для грошей :)")
		}
	case RejectAd:
		err = s.uc.UpdateAdChanEntry(models.AdvertisementChannel{
			ID:     uuid.FromStringOrNil(adChanID),
			Status: models.AdChanRejected,
		})
		if err != nil {
			zap.L().Error("failed to update advertisement status", zap.Error(err))

			return types.NewCustomMessageConfig(
				tgbotapi.NewMessage(respondTo, "Не вдалося відхилити рекламу"),
				nil,
				true,
				true,
				false,
			)
		} else {
			s.setState(respondTo, stateData{
				state:    StateWaitForRejectReason,
				adChanID: adChanID,
			})

			return types.NewCustomMessageConfig(
				tgbotapi.NewMessage(respondTo, "Поділіться, будь-ласка, причиною відмови:"),
				nil,
				false,
				false,
				false,
			)

		}
	}

	return types.NewCustomMessageConfig(
		msg,
		nil,
		true,
		false,
		false,
	)
}

func (s *Service) saveRejectionReason(respondTo int64, adChanID, reason string) types.CustomMessage {
	s.resetState(respondTo)
	err := s.uc.UpdateAdChanEntry(models.AdvertisementChannel{
		ID:              uuid.FromStringOrNil(adChanID),
		RejectionReason: reason,
	})
	if err != nil {
		zap.L().Error("failed to update advertisement status", zap.Error(err))
	}

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, "Дякую! Ми врахуємо це."),
		nil,
		true,
		false,
		false,
	)
}
