package transport

import (
	"advertiser/shared/pkg/service/repo/models"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

func (s *Transport) RunNotificationService() {
	go func() {
		for {
			s.NotifyAboutNewAds(models.AdChanCreated)
			time.Sleep(30 * time.Second)
		}
	}()
}

func (s *Transport) NotifyAboutNewAds(status models.AdChanStatus) {
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

func (s *Transport) PostAdvertisement(advertisementChannelID string) error {
	advertisementChannel, err := s.uc.GetAdChanDetails(advertisementChannelID)
	if err != nil {
		zap.L().Error("failed to get ad details", zap.Error(err))
		return nil
	}

	postedMsg, err := s.tgBotApi.Send(tgbotapi.NewMessage(
		advertisementChannel.ChannelID,
		fmt.Sprintf("%s", advertisementChannel.AdMessage)),
	)
	if err != nil {
		zap.L().Error("failed to post advertisement", zap.Error(err))
		return err
	}

	err = s.uc.UpdateAdChanStatus(advertisementChannelID, models.AdChanPosted)
	if err != nil {
		zap.L().Error("failed to update advertisement status", zap.Error(err))
	}

	err = s.uc.SetAdChanMessageID(advertisementChannelID, postedMsg.MessageID)
	if err != nil {
		zap.L().Error("failed to set message id", zap.Error(err))
	}

	return nil
}
