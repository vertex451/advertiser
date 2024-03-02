package transport

import (
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/pkg/service/transport"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func (t *Transport) allTopicsWithCoverage(respondTo int64) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	topics, err := t.uc.AllTopicsWithCoverage()
	if err != nil {
		zap.L().Error("failed to list topics", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, "failed to list topics")
	} else {
		var res []string
		for topic, coverage := range topics {
			if coverage > 0 {
				res = append(res, fmt.Sprintf("%s: %v subscribers", topic, coverage))
			}
		}
		msg = tgbotapi.NewMessage(respondTo, strings.Join(res, "\n"))
	}

	msg = transport.AddNavigationButtons(msg, nil)

	return &msg
}

func (t *Transport) createCampaign(respondTo int64, campaignName string) *tgbotapi.MessageConfig {
	t.resetState(respondTo)

	var msg tgbotapi.MessageConfig
	campaignID, err := t.uc.CreateCampaign(respondTo, campaignName)
	if err != nil {
		zap.L().Error("failed to create campaign", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to create campaign. Error: %v", err))
		msg = transport.AddNavigationButtons(msg, nil)
	} else {
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Campaing %s created!", campaignName))
		msg = transport.AddNavigationButtons(msg, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Create my first Ad in %s", campaignName),
				fmt.Sprintf("%s/%v", CreateAd, campaignID)),
		))
	}

	return &msg
}

func (t *Transport) listMyCampaigns(respondTo int64) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	myCampaigns, err := t.uc.ListMyCampaigns(respondTo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg = tgbotapi.NewMessage(respondTo, "You don't have campaigns")
		} else {
			zap.L().Error("failed to list campaigns", zap.Error(err))
			msg = tgbotapi.NewMessage(respondTo, "failed to list campaigns")
		}

		msg = transport.AddNavigationButtons(msg, nil)

		return &msg
	}

	var buttons []tgbotapi.InlineKeyboardButton
	if len(myCampaigns) > 0 {
		for _, campaign := range myCampaigns {
			buttons = append(buttons,
				tgbotapi.NewInlineKeyboardButtonData(
					campaign.Name, fmt.Sprintf("%s/%s", CampaignDetails, campaign.ID)),
			)
		}
		msg = tgbotapi.NewMessage(respondTo, "Select a campaign:")
	} else {
		msg = tgbotapi.NewMessage(respondTo, "You don't have campaigns")
		buttons = append(buttons,
			tgbotapi.NewInlineKeyboardButtonData(
				"Create my first campaign", CreateCampaign),
		)
	}

	msg = transport.AddNavigationButtons(msg, buttons)

	return &msg
}

func (t *Transport) campaignDetails(respondTo int64, rawCampaignID string) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig

	campaignID, err := uuid.FromString(rawCampaignID)
	if err != nil {
		zap.L().Error("failed to parse campaignID into uuid", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to read campaignID. Error: %v", err))
		msg = transport.AddNavigationButtons(msg, nil)

		return &msg
	}

	campaignDetails, err := t.uc.CampaignDetails(campaignID)
	if err != nil {
		zap.L().Error("failed to get campaignDetails", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to create an Ad. Error: %v", err))
		msg = transport.AddNavigationButtons(msg, nil)

		return &msg
	}

	var buttons []tgbotapi.InlineKeyboardButton
	if len(campaignDetails.Advertisements) > 0 {
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("%s advertisements:", campaignDetails.Name))
		for _, ad := range campaignDetails.Advertisements {
			buttons = append(buttons,
				tgbotapi.NewInlineKeyboardButtonData(
					ad.Name, fmt.Sprintf("%s/%s", AdDetails, ad.ID)),
			)
		}
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("%s advertisements", campaignDetails.Name))
	} else {
		msg = tgbotapi.NewMessage(respondTo, "You don't have advertisements yet")
	}

	buttons = append(buttons,
		tgbotapi.NewInlineKeyboardButtonData(
			"Create an advertisement", fmt.Sprintf("%s/%s", CreateAd, campaignID)),
	)

	msg = transport.AddNavigationButtons(msg, buttons)

	return &msg
}

func (t *Transport) upsertAd(respondTo int64, campaignID, adID, input string) *tgbotapi.MessageConfig {
	t.resetState(respondTo)

	var msg tgbotapi.MessageConfig
	ad, err := parseAndValidateCreateAdInput(campaignID, adID, input)
	if err != nil {
		zap.L().Error("failed to parse an input", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to create an Ad. Error: %v", err))
		msg = transport.AddNavigationButtons(msg, nil)

		return &msg
	}

	_, err = t.uc.UpsertAd(*ad)
	if err != nil {
		zap.L().Error("failed to create an ad", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to create an Ad. Error: %v", err))
	} else {
		if ad.ID == uuid.Nil {
			msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Advertisement %s created!", ad.Name))
		} else {
			msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Advertisement %s updated!!", ad.Name))
		}
	}

	msg = transport.AddNavigationButtons(msg, nil)

	return &msg
}

func (t *Transport) GetAdDetails(respondTo int64, rawID string) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig

	id, err := uuid.FromString(rawID)
	if err != nil {
		zap.L().Error("failed to parse id", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to parse advertisement id. Error: %v", err))
		msg = transport.AddNavigationButtons(msg, nil)

		return &msg
	}

	var buttons []tgbotapi.InlineKeyboardButton

	ad, err := t.uc.GetAdDetails(id)
	if err != nil {
		zap.L().Error("failed to get advertisement details", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to get advertisement details. Error: %v", err))
	} else {
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf(`
Name: %s
TargetTopics: %s
BudgetUSD: %v
Message: %s

Estimated coverage: %v
`,
			ad.Name,
			strings.Join(ad.Topics, ", "),
			ad.Budget,
			ad.Message,
			ad.Coverage,
		))

		if ad.Status == models.AdsStatusCreated {
			buttons = append(buttons,
				tgbotapi.NewInlineKeyboardButtonData(
					"Edit", fmt.Sprintf("%s/%s", EditAd, ad.ID)),
				tgbotapi.NewInlineKeyboardButtonData(
					"Run", fmt.Sprintf("%s/%s", RunAd, ad.ID)),
			)
		} else {
			buttons = append(buttons,
				tgbotapi.NewInlineKeyboardButtonData(
					"Pause", fmt.Sprintf("%s/%s", PauseAd, ad.ID)),
				tgbotapi.NewInlineKeyboardButtonData(
					"Finish", fmt.Sprintf("%s/%s", FinishAd, ad.ID)),
			)
		}
	}

	msg = transport.AddNavigationButtons(msg, buttons)

	return &msg

}

func (t *Transport) RunAd(respondTo int64, rawID string) *tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig

	id, err := uuid.FromString(rawID)
	if err != nil {
		zap.L().Error("failed to parse id", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to parse advertisement id. Error: %v", err))
		msg = transport.AddNavigationButtons(msg, nil)

		return &msg
	}

	var buttons []tgbotapi.InlineKeyboardButton
	err = t.uc.RunAd(id)
	if err != nil {
		zap.L().Error("failed ailed to run advertisement", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to run advertisement. Error: %v", err))
	} else {
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Advertising is running! It will start appearing in channels after an approval from channel owners"))
	}

	msg = transport.AddNavigationButtons(msg, buttons)

	return &msg

}

func parseAndValidateCreateAdInput(rawCampaignID, rawAdID, rawInput string) (*models.Advertisement, error) {
	requiredFields := []string{"Name", "TargetTopics", "BudgetUSD", "CostPerView", "Message"}

	params := parseValues(rawInput)
	for _, field := range requiredFields {
		if val, ok := params[field]; !ok || val == "" {
			return nil, errors.New("missing required field: " + field)
		}
	}

	var ad models.Advertisement
	if rawCampaignID != "" {
		var campaignID, err = uuid.FromString(rawCampaignID)
		if err != nil {
			return nil, errors.New("failed to parse campaignID into uuid")
		}
		ad.CampaignID = campaignID
	}

	if rawAdID != "" {
		var adID, err = uuid.FromString(rawAdID)
		if err != nil {
			return nil, errors.New("failed to parse adID into uuid")
		}
		ad.ID = adID
	}

	for key, value := range params {
		switch key {
		case "Name":
			ad.Name = value
		case "BudgetUSD":
			budget := 0
			_, err := fmt.Sscanf(value, "%d", &budget)
			if err != nil {
				return nil, errors.New("invalid budget format")
			}
			if budget <= 0 {
				return nil, errors.New("budget should be greater than 0")
			}
			ad.Budget = budget
		case "CostPerView":
			costPerView, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return nil, errors.New("invalid budget format")
			}
			if costPerView <= 0 {
				return nil, errors.New("cost per view should be greater than 0")
			}
			ad.CostPerView = float32(costPerView)
		case "Message":
			ad.Message = value
		case "TargetTopics":
			topics := strings.Split(value, ",")
			for _, topicName := range topics {
				ad.TargetTopics = append(ad.TargetTopics, models.Topic{ID: strings.TrimSpace(topicName)})
			}
		}
	}

	return &ad, nil
}

func parseValues(input string) map[string]string {
	values := make(map[string]string)

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		values[key] = value
	}

	return values
}
