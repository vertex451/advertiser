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
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const AdMessageMaxLength = 700

func (t *Transport) allTopicsWithCoverage(respondTo int64) types.CustomMessage {
	var msg tgbotapi.MessageConfig
	topics, err := t.uc.AllTopicsWithCoverage()
	if err != nil {
		zap.L().Error("failed to list topics", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, "failed to list topics")
	} else {
		var res []string
		for _, topic := range topics {
			res = append(res, fmt.Sprintf("%s: %v subscribers", topic.Name, topic.Coverage))
		}
		msg = tgbotapi.NewMessage(respondTo, strings.Join(res, "\n"))
	}

	return types.NewCustomMessageConfig(
		msg,
		nil,
		true,
		false,
	)
}

func (t *Transport) createCampaign(respondTo int64, campaignName string) types.CustomMessage {
	t.resetState(respondTo)

	campaignID, err := t.uc.CreateCampaign(respondTo, campaignName)
	if err != nil {
		zap.L().Error("failed to create campaign", zap.Error(err))

		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to create campaign. Error: %v", err)),
			nil,
			true,
			false,
		)
	}

	return types.NewCustomMessageConfig(
		tgbotapi.NewMessage(respondTo, fmt.Sprintf("Campaign %s created!", campaignName)),
		[][]tgbotapi.InlineKeyboardButton{{
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Create my first Ad in %s", campaignName),
				fmt.Sprintf("%s/%v", CreateAd, campaignID)),
		}},
		true,
		false,
	)
}

func (t *Transport) listMyCampaigns(respondTo int64) types.CustomMessage {
	var msg tgbotapi.MessageConfig
	myCampaigns, err := t.uc.ListMyCampaigns(respondTo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg = tgbotapi.NewMessage(respondTo, "You don't have campaigns")
		} else {
			zap.L().Error("failed to list campaigns", zap.Error(err))
			msg = tgbotapi.NewMessage(respondTo, "failed to list campaigns")
		}

		return types.NewCustomMessageConfig(
			msg,
			nil,
			true,
			false,
		)
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

	return types.NewCustomMessageConfig(
		msg,
		transport.MakeTwoButtonsInARow(buttons),
		true,
		false,
	)
}

func (t *Transport) campaignDetails(respondTo int64, rawCampaignID string) types.CustomMessage {
	var msg tgbotapi.MessageConfig

	campaignID, err := uuid.FromString(rawCampaignID)
	if err != nil {
		zap.L().Error("failed to parse campaignID into uuid", zap.Error(err))

		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to read campaignID. Error: %v", err)),
			nil,
			true,
			false,
		)
	}

	campaignDetails, err := t.uc.CampaignDetails(campaignID)
	if err != nil {
		zap.L().Error("failed to get campaignDetails", zap.Error(err))

		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to get an campaignDetails. Error: %v", err)),
			nil,
			true,
			false,
		)
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
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("%s advertisements:", campaignDetails.Name))
	} else {
		msg = tgbotapi.NewMessage(respondTo, "You don't have advertisements yet")
	}

	rows := transport.MakeTwoButtonsInARow(buttons)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
		"Create an advertisement", fmt.Sprintf("%s/%s", CreateAd, campaignID))))

	return types.NewCustomMessageConfig(
		msg,
		rows,
		true,
		false,
	)
}

func (t *Transport) createAdMetadata(respondTo int64, campaignID, input string) types.CustomMessage {
	t.resetState(respondTo)

	ad, err := parseAndValidateCreateAdInput(campaignID, input)
	if err != nil {
		zap.L().Error("failed to parse an input", zap.Error(err))

		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to create an Ad. Error: %v", err)),
			nil,
			true,
			false,
		)
	}

	err = t.uc.UpsertAd(ad)
	if err != nil {
		zap.L().Error("failed to create an ad", zap.Error(err))
		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to create an Ad. Error: %v", err)),
			nil,
			true,
			false,
		)
	}

	t.setState(respondTo, stateData{
		state: StateCreateAdMessage,
		adID:  ad.ID.String(),
	})

	photoMsg := tgbotapi.NewPhoto(
		respondTo,
		tgbotapi.FilePath("/Users/vertex451/workplace/silverspase/tg-bot-images/example.jpg"),
	)
	photoMsg.Caption = `<b>Now send an advertisement message!</b>

You can create rich messages using the following formatting options:
<b>bold</b>, <i>italic</i>, <strike>strikethrough</strike>, <u>underline</u>,  <a href="https://www.google.com/">link</a>, <code>monospace</code> or even <tg-spoiler>spoiler</tg-spoiler>!
You can attach the image as well!

Just send me a message as you usually do and I will keep the formatting.`
	photoMsg.ParseMode = tgbotapi.ModeHTML

	return types.NewCustomPhotoConfig(
		photoMsg,
		nil,
		false,
		false,
	)
}

func (t *Transport) createAdMessage(respondTo int64, adID, text string, inputEntities []tgbotapi.MessageEntity, photos []tgbotapi.PhotoSize) types.CustomMessage {
	t.resetState(respondTo)

	var err error
	var msg tgbotapi.MessageConfig

	if len(text) > AdMessageMaxLength {
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf(`
Message is too long. Max length is %d.
Please, send new message with shorter length:
`, AdMessageMaxLength))

		t.setState(respondTo, stateData{
			state: StateCreateAdMessage,
			adID:  adID,
		})

		return types.NewCustomMessageConfig(
			msg,
			nil,
			false,
			true,
		)
	}

	messageText := transport.AddFooter(text)

	var photoPath string

	photoPath, err = t.saveMsgPhoto(photos)
	if err != nil {
		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to save photo. Error: %v", err)),
			nil,
			true,
			false,
		)
	}

	var entities []models.MsgEntity
	for _, entity := range inputEntities {
		entities = append(entities, models.MsgEntity{
			Type:     entity.Type,
			Offset:   entity.Offset,
			Length:   entity.Length,
			URL:      entity.URL,
			Language: entity.Language,
		})
	}

	err = t.uc.UpsertAd(&models.Advertisement{
		ID:          uuid.FromStringOrNil(adID),
		MsgText:     messageText,
		MsgImageURL: photoPath,
		MsgEntities: entities,
	})
	if err != nil {
		zap.L().Error("failed to create an advertisement adMessage", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to create an Advertisement. Error: %v", err))
	} else {
		msg = tgbotapi.NewMessage(respondTo, "Advertisement created!")
	}

	return types.NewCustomMessageConfig(
		msg,
		[][]tgbotapi.InlineKeyboardButton{{
			tgbotapi.NewInlineKeyboardButtonData(
				"View message", fmt.Sprintf("%s/%s", constants.ViewAdMessage, adID)),
		}},
		true,
		true,
	)
}

func (t *Transport) saveMsgPhoto(photos []tgbotapi.PhotoSize) (string, error) {
	if len(photos) < 1 {
		return "", nil
	}

	saveDir := "/Users/vertex451/workplace/silverspase/tg-bot-images"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		zap.L().Fatal("Error creating directory", zap.Error(err))

		return "", errors.Errorf("Error creating directory: %v", err)
	}

	// Get the photo
	photo := photos[len(photos)-1]

	file, err := t.tgBotApi.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
	if err != nil {
		zap.L().Error("Error getting file", zap.Error(err))

		return "", errors.Errorf("Error getting file: %v", err)
	}
	fileID := file.FileID

	// Get the direct URL of the photo
	fileURL, err := t.tgBotApi.GetFileDirectURL(fileID)
	if err != nil {
		zap.L().Error("Error getting file URL", zap.Error(err))

		return "", errors.Errorf("Error getting file URL: %v", err)
	}

	// Download the photo
	resp, err := http.Get(fileURL)
	if err != nil {
		zap.L().Error("Error downloading file", zap.Error(err))

		return "", errors.Errorf("Error downloading file: %v", err)
	}
	defer resp.Body.Close()

	// Create the file to save the photo
	savePath := filepath.Join(saveDir, filepath.Base(fileID)+filepath.Ext(fileURL))
	saveFile, err := os.Create(savePath)
	if err != nil {
		zap.L().Error("Error creating file", zap.Error(err))

		return "", errors.Errorf("Error creating file: %v", err)
	}
	defer saveFile.Close()

	// Save the photo data to the file
	_, err = io.Copy(saveFile, resp.Body)
	if err != nil {
		zap.L().Error("Error saving file", zap.Error(err))

		return "", errors.Errorf("Error saving file: %v", err)
	}

	return savePath, nil
}

func (t *Transport) getAdDetails(respondTo int64, rawID string) types.CustomMessage {
	var msg tgbotapi.MessageConfig

	id, err := uuid.FromString(rawID)
	if err != nil {
		zap.L().Error("failed to parse id", zap.Error(err))

		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to parse advertisement id. Error: %v", err)),
			nil,
			true,
			false,
		)
	}

	var buttons []tgbotapi.InlineKeyboardButton

	ad, err := t.uc.GetAdDetails(id)
	if err != nil {
		zap.L().Error("failed to get advertisement details", zap.Error(err))
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf("Failed to get advertisement details. Error: %v", err))
	} else {
		msg = tgbotapi.NewMessage(respondTo, fmt.Sprintf(`
Name: %s
Status: %s
TargetTopics: %s
BudgetUSD: %v
`,
			ad.Name,
			ad.Status,
			ad.GetTopics(),
			ad.Budget,
		))

		//if ad.Status == models.AdsStatusCreated {
		buttons = append(buttons,
			tgbotapi.NewInlineKeyboardButtonData(
				"View message", fmt.Sprintf("%s/%s", constants.ViewAdMessage, ad.ID)),
		)
		//} else {
		//buttons = append(buttons,
		//	tgbotapi.NewInlineKeyboardButtonData(
		//		"Pause", fmt.Sprintf("%s/%s", PauseAd, ad.ID)),
		//	tgbotapi.NewInlineKeyboardButtonData(
		//		"Finish", fmt.Sprintf("%s/%s", FinishAd, ad.ID)),
		//)
		//}
	}

	return types.NewCustomMessageConfig(
		msg,
		[][]tgbotapi.InlineKeyboardButton{buttons},
		true,
		true,
	)
}

func (t *Transport) viewAdMessage(userID int64, variable string) types.CustomMessage {
	ad, err := t.uc.GetAdDetails(uuid.FromStringOrNil(variable))
	if err != nil {
		zap.L().Error("failed to get advertisement details", zap.Error(err))

		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(userID, fmt.Sprintf("Failed to get advertisement details. Error: %v", err)),
			nil,
			true,
			false,
		)
	}

	return transport.ComposeAdMessage(
		userID,
		*ad,
		[][]tgbotapi.InlineKeyboardButton{
			{tgbotapi.NewInlineKeyboardButtonData(
				"Run", fmt.Sprintf("%s/%s", RunAd, ad.ID)),
				tgbotapi.NewInlineKeyboardButtonData(
					"Delete", fmt.Sprintf("%s/%s", DeleteAd, ad.ID)),
			},
		},
		true,
		false,
	)
}

func (t *Transport) RunAd(userID int64, rawID string) types.CustomMessage {
	var msg tgbotapi.MessageConfig

	id, err := uuid.FromString(rawID)
	if err != nil {
		zap.L().Error("failed to parse id", zap.Error(err))

		return types.NewCustomMessageConfig(
			tgbotapi.NewMessage(userID, fmt.Sprintf("Failed to parse advertisement id. Error: %v", err)),
			nil,
			true,
			false,
		)
	}

	err = t.uc.RunAd(id)
	if err != nil {
		zap.L().Error("failed ailed to run advertisement", zap.Error(err))
		msg = tgbotapi.NewMessage(userID, fmt.Sprintf("Failed to run advertisement. Error: %v", err))
	} else {
		msg = tgbotapi.NewMessage(userID, fmt.Sprintf("Advertising is running! It will start appearing in channels after an approval from channel owners"))
	}

	return types.NewCustomMessageConfig(
		msg,
		nil,
		true,
		false,
	)
}

func parseAndValidateCreateAdInput(rawCampaignID, rawInput string) (*models.Advertisement, error) {
	requiredFields := []string{"Name", "TargetTopics", "BudgetUSD", "CostPerView"}

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
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(strings.Join(parts[1:], ":"))
		values[key] = value
	}

	return values
}
