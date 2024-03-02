package models

type ChannelAdmin struct {
	ChannelID int64 `gorm:"primaryKey"`
	UserID    int64 `gorm:"primaryKey"`
	Role      string
}
