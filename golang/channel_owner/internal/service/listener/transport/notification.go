package transport

import (
	"advertiser/shared/pkg/service/repo/models"
	"advertiser/shared/pkg/service/transport"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

func (s *Service) NotifyChannelOwnersAboutNewAds(status models.AdChanStatus) {
	res, err := s.uc.GetAdsChannelByStatus(status)
	if err != nil {
		return
	}

	//var firstNotificationSent []uuid.UUID
	var msg tgbotapi.MessageConfig
	for userID, adsChannels := range res {
		msg = tgbotapi.NewMessage(userID, fmt.Sprintf(`
You have %d advertisements to moderate.
Click on /%s to view them.
`, len(adsChannels), Moderate))
		_, err = s.tgBotApi.Send(msg)
		if err != nil {
			zap.L().Error("failed to send message to moderation", zap.Error(err))
			continue
		}
	}
}

func (s *Service) PostAdvertisement(adChanID string) error {
	adChan, err := s.uc.GetAdChanDetails(adChanID)
	if err != nil {
		zap.L().Error("failed to get ad details", zap.Error(err))
		return nil
	}

	msg := transport.ComposeAdMessage(
		adChan.Channel.ID,
		adChan.Advertisement,
		nil,
		false,
		false,
	)
	postedMsg, err := s.tgBotApi.Send(msg)
	if err != nil {
		zap.L().Error("failed to post advertisement", zap.Error(err))
		return err
	}

	deleteAt := time.Now().Add(5 * time.Second)
	err = s.ScheduleMsgDeletionAtTime(
		adChanID,
		adChan.ChannelID,
		postedMsg.MessageID,
		deleteAt,
	)
	if err != nil {
		zap.L().Error("failed to schedule message deletion", zap.Error(err))
		return err
	}

	err = s.uc.UpdateAdChanEntry(models.AdvertisementChannel{
		ID:                  uuid.FromStringOrNil(adChanID),
		Status:              models.AdChanPosted,
		ChannelPostID:       postedMsg.MessageID,
		DeletionScheduledAt: deleteAt,
	})
	if err != nil {
		zap.L().Error("failed to set message id", zap.Error(err))
	}

	return nil
}

func (s *Service) DeleteAdvertisement(adChanID string, channelID int64, messageID int) error {
	editMessageConfig := tgbotapi.NewEditMessageText(
		channelID,
		messageID,
		transport.AddFooter("Advertisement is finished."),
	)

	_, err := s.tgBotApi.Send(editMessageConfig)
	if err != nil {
		zap.L().Error("failed to delete advertisement", zap.Error(err))
		return err
	}

	err = s.uc.UpdateAdChanEntry(models.AdvertisementChannel{
		ID:     uuid.FromStringOrNil(adChanID),
		Status: models.AdChanFinished,
	})
	if err != nil {
		zap.L().Error("failed to set message id", zap.Error(err))
		return err
	}

	return nil
}
