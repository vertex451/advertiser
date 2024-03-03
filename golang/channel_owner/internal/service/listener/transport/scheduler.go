package transport

import (
	"advertiser/shared/pkg/service/repo/models"
	"fmt"
	"go.uber.org/zap"
	"time"
)

const EveryMinute = "*/1 * * * *"

func (s *Transport) RunNotificationService() {
	s.StartChannelOwnerNewAdsNotificator()
	s.StartNewAdsChecker()
}

func (s *Transport) StartNewAdsChecker() {
	_, err := s.cron.AddFunc(EveryMinute, func() {
		s.uc.CheckForNewAds()
	})
	if err != nil {
		zap.L().Error("failed to run new ads checker", zap.Error(err))
	}
}

func (s *Transport) StartChannelOwnerNewAdsNotificator() {
	_, err := s.cron.AddFunc(EveryMinute, func() {
		s.NotifyChannelOwnersAboutNewAds(models.AdChanCreated)
	})
	if err != nil {
		zap.L().Error("failed to run new ads notification service", zap.Error(err))
	}
}

func (s *Transport) ScheduleMsgDeletion(adChanID string, channelID int64, messageID int, at time.Time) error {
	_, err := s.cron.AddFunc(
		fmt.Sprintf("%d %d %d %d *", at.Minute()+1, at.Hour(), at.Day(), int(at.Month())),
		func() {
			s.DeleteAdvertisement(adChanID, channelID, messageID)
		})

	return err
}
