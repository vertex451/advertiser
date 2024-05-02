package postgresql

import (
	"advertiser/shared/pkg/service/constants"
	"advertiser/shared/pkg/service/repo/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	uuid "github.com/satori/go.uuid"
)

func (r *Repository) AllTopics() (res []string, err error) {
	var allTopics []*models.Topic
	if err = r.Db.Find(&allTopics).Error; err != nil {
		return nil, err
	}

	for _, topic := range allTopics {
		res = append(res, topic.ID)
	}

	return res, nil
}

func (r *Repository) StoreInitialChannelData(admins []tgbotapi.ChatMember, channel models.Channel) error {
	if result := r.Db.Save(&channel); result.Error != nil {
		return result.Error
	}

	if err := r.upsertAdminAndCreateLinks(channel.ID, admins); err != nil {
		return err
	}

	return nil
}

func (r *Repository) upsertAdminAndCreateLinks(chatID int64, admins []tgbotapi.ChatMember) (err error) {
	for _, oneAdmin := range admins {
		var admin models.User
		if tx := r.Db.FirstOrCreate(&admin, models.User{
			ID:     oneAdmin.User.ID,
			Handle: oneAdmin.User.UserName,
		}); tx.Error != nil {
			return err
		}

		if tx := r.Db.Save(&models.ChannelAdmin{
			ChannelID: chatID,
			UserID:    admin.ID,
			Role:      oneAdmin.Status,
		}); tx.Error != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) ListMyChannels(userID int64) (map[int64]string, error) {
	var user models.User
	if err := r.Db.Model(&models.User{}).Preload("Channels").First(&user, userID).Error; err != nil {
		return nil, err
	}

	result := make(map[int64]string)
	for _, channel := range user.Channels {
		result[channel.ID] = channel.Title
	}

	return result, nil
}

func (r *Repository) GetChannelInfo(channelID int64) (*models.Channel, error) {
	var channel models.Channel
	if err := r.Db.Preload("Topics").First(&channel, channelID).Error; err != nil {
		return nil, err
	}

	return &channel, nil
}

func (r *Repository) UpdateChannelTopics(channelID int64, newTopicNames []string) (err error) {
	// Find the channel by ID
	var channel models.Channel
	if err = r.Db.Preload("Topics").First(&channel, channelID).Error; err != nil {
		return err
	}

	// Clear existing topics associated with the channel
	if err = r.Db.Model(&channel).Association("Topics").Clear(); err != nil {
		return err
	}

	// Fetch new topics from the database
	var newTopics []models.Topic
	if err = r.Db.Where("id IN ?", newTopicNames).Find(&newTopics).Error; err != nil {
		return err
	}

	// Check if all provided topics were found
	if len(newTopics) != len(newTopicNames) {
		missingTopics := findMissingTopics(newTopics, newTopicNames)
		return fmt.Errorf("topics not found: %v", missingTopics)
	}

	// Assign new topics to the channel
	channel.Topics = newTopics

	return r.Db.Save(&channel).Error
}

func (r *Repository) UpdateChannelCostPerMile(channelID int64, costPerMile float64) error {
	return r.Db.Model(&models.Channel{}).Where("id = ?", channelID).Update("cost_per_mile", costPerMile).Error
}

func (r *Repository) UpdateChannelLocation(channelID int64, location constants.Location) error {
	return r.Db.Model(&models.Channel{}).Where("id = ?", channelID).Update("location", location).Error
}

func (r *Repository) DeleteChannel(chatID int64) error {
	err := r.Db.Delete(&models.Channel{}, chatID).Error
	if err != nil {
		return err
	}

	return r.DeleteChannelAdminsEntries(chatID)
}

func (r *Repository) DeleteChannelAdminsEntries(chatID int64) error {
	return r.Db.Delete(&models.ChannelAdmin{}, "channel_id = ?", chatID).Error
}

func (r *Repository) GetAdsToModerateByUserID(id int64) ([]models.AdvertisementChannel, error) {
	var ads []models.AdvertisementChannel
	// If you need to support each admin, remove filter by role
	if err := r.Db.
		Preload("Advertisement").
		Preload("Channel.ChannelAdmins", "user_id = ? AND role = ?", id, constants.StatusCreator).
		Where("status = ?", models.AdsStatusCreated).Find(&ads).Error; err != nil {
		return nil, err
	}

	return ads, nil
}

func (r *Repository) GetAdChanDetails(id string) (*models.AdvertisementChannel, error) {
	var ad models.AdvertisementChannel
	if err := r.Db.
		Preload("Advertisement.MsgEntities").
		Preload("Channel").
		First(&ad, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &ad, nil
}

func (r *Repository) GetAdMessageByAdChanID(adChanID uuid.UUID) (*models.Advertisement, error) {
	//var advertisement models.Advertisement
	var adChan models.AdvertisementChannel
	if err := r.Db.
		Preload("Advertisement.MsgEntities").
		First(&adChan, adChanID).Error; err != nil {
		return nil, err
	}

	return &adChan.Advertisement, nil
}

func (r *Repository) ReportBug(userID int64, message string) error {
	return r.Db.Save(&models.BugReport{
		ReportedBy: userID,
		Message:    message,
	}).Error
}

func (r *Repository) RequestFeature(userID int64, message string) error {
	return r.Db.Save(&models.FeatureRequest{
		RequestedBy: userID,
		Message:     message,
	}).Error
}

func findMissingTopics(foundTopics []models.Topic, providedNames []string) []string {
	foundMap := make(map[string]bool)
	for _, topic := range foundTopics {
		foundMap[topic.ID] = true
	}

	var missingTopics []string
	for _, name := range providedNames {
		if !foundMap[name] {
			missingTopics = append(missingTopics, name)
		}
	}
	return missingTopics
}
