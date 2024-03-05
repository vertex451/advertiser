package transport

import (
	"advertiser/shared/pkg/service/repo/models"
	"fmt"
	"go.uber.org/zap"
	"time"
)

const EveryMinute = "*/1 * * * *"

func (s *Transport) RunNotificationService() {
	s.StartChannelOwnerNewAdsChecker()
	s.StartNewAdsChecker()
	s.StartOverspendingChecker()
}

func (s *Transport) StartNewAdsChecker() {
	_, err := s.cron.AddFunc(EveryMinute, func() {
		s.uc.CheckForNewAds()
	})
	if err != nil {
		zap.L().Error("failed to run new ads checker", zap.Error(err))
	}
}

func (s *Transport) StartChannelOwnerNewAdsChecker() {
	_, err := s.cron.AddFunc(EveryMinute, func() {
		s.NotifyChannelOwnersAboutNewAds(models.AdChanCreated)
	})
	if err != nil {
		zap.L().Error("failed to run new ads notification service", zap.Error(err))
	}
}

func (s *Transport) StartOverspendingChecker() {
	_, err := s.cron.AddFunc(EveryMinute, func() {
		s.PreventOverspending()
	})
	if err != nil {
		zap.L().Error("failed to prevent overspending", zap.Error(err))
	}
}

func (s *Transport) ScheduleMsgDeletionAtTime(adChanID string, channelID int64, messageID int, at time.Time) error {
	_, err := s.cron.AddFunc(
		fmt.Sprintf("%d %d %d %d *", at.Minute()+1, at.Hour(), at.Day(), int(at.Month())),
		func() {
			s.DeleteAdvertisement(adChanID, channelID, messageID)
		})

	return err
}
