package models

type ChannelTopic struct {
	ChannelID int64  `gorm:"primaryKey"`
	TopicID   string `gorm:"primaryKey"`
}
