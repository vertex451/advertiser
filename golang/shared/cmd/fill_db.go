package main

import (
	"advertiser/shared/config/config"
	"advertiser/shared/pkg/service/repo"
	"advertiser/shared/pkg/service/repo/models"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(cfg *config.Config) *Repository {
	db := repo.New(cfg)

	return &Repository{
		db: db,
	}
}

func main() {
	r := New(config.Load())
	r.FillTopics()
	r.fillChannelOwnerData()
	r.fillAgencyData()

	fmt.Println("Database filled")
}

func (r *Repository) fillChannelOwnerData() {
	r.fillChannels()
}

func (r *Repository) FillTopics() {
	topics := []models.Topic{{ID: "art"}, {ID: "books"}, {ID: "food"}, {ID: "pets"}, {ID: "sport"}}

	for _, topic := range topics {
		result := r.db.FirstOrCreate(&topic)
		if result.Error != nil {
			fmt.Printf("failed to create topic %s: %s\n", topic.ID, result.Error)
		}
	}
}

func (r *Repository) fillChannels() {
	channels := []models.Channel{
		{
			ID:          -1002134289719,
			Handle:      "pets132213",
			IsChannel:   true,
			Title:       "Pets channel",
			Subscribers: 2,
		},
		{
			ID:          -1002049183103,
			Handle:      "PublicBooksChannel451",
			IsChannel:   true,
			Title:       "BooksChannel",
			Subscribers: 3,
		},
	}

	for _, channel := range channels {
		result := r.db.FirstOrCreate(&channel)
		if result.Error != nil {
			fmt.Printf("failed to create channel %s: %s\n", channel.ID, result.Error)
		}
	}

	r.fillChannelTopics()
}

func (r *Repository) fillChannelTopics() {
	channelTopics := []models.ChannelTopic{
		{
			ChannelID: -1002049183103,
			TopicID:   "books",
		},
		{
			ChannelID: -1002049183103,
			TopicID:   "food",
		},
		{
			ChannelID: -1002134289719,
			TopicID:   "food",
		},
		{
			ChannelID: -1002134289719,
			TopicID:   "pets",
		},
	}

	for _, channelTopic := range channelTopics {
		result := r.db.FirstOrCreate(&channelTopic)
		if result.Error != nil {
			fmt.Printf("failed to create channel topic %v: %s\n", channelTopic.ChannelID, result.Error)
		}
	}

	r.fillUsers()
}

func (r *Repository) fillUsers() {
	users := []models.User{
		{
			ID:     6761224677,
			Handle: "channel_monetizer_bot",
		},
		{
			ID:     6406834985,
			Handle: "skydreamer451",
		},
	}

	for _, user := range users {
		result := r.db.FirstOrCreate(&user)
		if result.Error != nil {
			fmt.Printf("failed to create user %v: %s\n", user.ID, result.Error)
		}
	}

	r.fillChannelAdmins()
}

func (r *Repository) fillChannelAdmins() {
	admins := []models.ChannelAdmin{
		{
			ChannelID: -1002134289719,
			UserID:    6761224677,
			Role:      "administrator",
		},
		{
			ChannelID: -1002134289719,
			UserID:    6406834985,
			Role:      "creator",
		},
		{
			ChannelID: -1002049183103,
			UserID:    6761224677,
			Role:      "administrator",
		},
		{
			ChannelID: -1002049183103,
			UserID:    6406834985,
			Role:      "creator",
		},
	}

	for _, admin := range admins {
		result := r.db.FirstOrCreate(&admin)
		if result.Error != nil {
			fmt.Printf("failed to create channel admin %s: %s\n", admin.ChannelID, result.Error)
		}
	}
}

func (r *Repository) fillAgencyData() {
	r.fillCampaigns()
}

func (r *Repository) fillCampaigns() {
	campaigns := []models.Campaign{
		{
			ID:     uuid.FromStringOrNil("1f97e147-95cd-46c3-b2e1-a2f750a486e8"),
			UserID: 6406834985,
			Name:   "Food",
		},
	}

	for _, campaign := range campaigns {
		result := r.db.FirstOrCreate(&campaign)
		if result.Error != nil {
			fmt.Printf("failed to create campaign %s: %s\n", campaign.ID, result.Error)
		}
	}

	r.fillAds()
}

func (r *Repository) fillAds() {
	ads := []models.Advertisement{
		{
			ID:          uuid.FromStringOrNil("25f9451e-1f65-426f-85ff-7735bc39fc41"),
			CampaignID:  uuid.FromStringOrNil("1f97e147-95cd-46c3-b2e1-a2f750a486e8"),
			Name:        "McDonald's",
			Status:      "created",
			Budget:      100,
			CostPerView: 0.1,
			MsgText:     "Try our new burger!",
		},
	}

	for _, ad := range ads {
		result := r.db.FirstOrCreate(&ad)
		if result.Error != nil {
			fmt.Printf("failed to create ad %s: %s\n", ad.ID, result.Error)
		}
	}

	r.fillAdTopics()
}

func (r *Repository) fillAdTopics() {
	adTopics := []models.AdvertisementTopic{
		{
			AdvertisementID: uuid.FromStringOrNil("25f9451e-1f65-426f-85ff-7735bc39fc41"),
			TopicID:         "food",
		},
	}

	for _, adTopic := range adTopics {
		result := r.db.FirstOrCreate(&adTopic)
		if result.Error != nil {
			fmt.Printf("failed to create ad topic %s: %s\n", adTopic.AdvertisementID, result.Error)
		}
	}
}
