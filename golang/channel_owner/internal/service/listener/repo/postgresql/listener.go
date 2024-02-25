package postgresql

import (
	"advertiser/shared/pkg/service/repo/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	if _, ok := r.channelIdByHandle.Load(channel.Handle); ok {
		return nil
	}

	if result := r.Db.Create(&channel); result.Error != nil {
		return result.Error
	}

	if err := r.upsertAdminAndCreateLinks(channel.ID, admins); err != nil {
		return err
	}

	r.channelIdByHandle.Store(channel.Handle, channel.ID)

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

		if tx := r.Db.Create(&models.ChannelAdmin{
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
	var newTopics []*models.Topic
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

func (r *Repository) DeleteChannel(chatID int64) error {
	return r.Db.Delete(&models.Channel{}, chatID).Error
}

func findMissingTopics(foundTopics []*models.Topic, providedNames []string) []string {
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
