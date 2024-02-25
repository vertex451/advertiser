package models

type AdsChannelStatus string

const (
	AdsChannelPendingApproval AdsChannelStatus = "pending_approval"
	AdsChannelApproved        AdsChannelStatus = "approved"
	AdsChannelRejected        AdsChannelStatus = "rejected"
	AdsChannelPosted          AdsChannelStatus = "posted"
	AdsChannelFinished        AdsChannelStatus = "finished"
)

type AdvertisementChannel struct {
	AdvertisementID int64
	ChannelID       int64
	Status          string
}
