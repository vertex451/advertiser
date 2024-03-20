package transport

import (
	"advertiser/shared/config/config"
	"advertiser/shared/pkg/service/repo/models"
	"fmt"
	"go.uber.org/zap"
	"time"
)

const Every2Seconds = "*/2 * * * * *"
const Every10Seconds = "*/10 * * * * *"
const EveryHour = "0 0 * * * *"

func getInterval(env string) string {
	switch env {
	case config.EnvIntegration:
		return Every2Seconds
	case config.EnvDev:
		return Every10Seconds
	case config.EnvProd:
		return EveryHour
	default:
		return Every10Seconds
	}
}

func (s *Service) RunNotificationService() {
	zap.L().Info("running notification service", zap.String("pattern", getInterval(s.env)))
	s.StartChannelOwnerNewAdsChecker()
	s.StartNewAdsChecker()
	s.StartOverspendingChecker()
}

func (s *Service) StartNewAdsChecker() {
	_, err := s.cron.AddFunc(getInterval(s.env), func() {
		s.uc.CheckForNewAds()
	})
	if err != nil {
		zap.L().Error("failed to run new ads checker", zap.Error(err))
	}
}

func (s *Service) StartChannelOwnerNewAdsChecker() {
	_, err := s.cron.AddFunc(getInterval(s.env), func() {
		s.NotifyChannelOwnersAboutNewAds(models.AdChanCreated)
	})
	if err != nil {
		zap.L().Error("failed to run new ads notification service", zap.Error(err))
	}
}

func (s *Service) StartOverspendingChecker() {
	_, err := s.cron.AddFunc(getInterval(s.env), func() {
		s.PreventOverspending()
	})
	if err != nil {
		zap.L().Error("failed to prevent overspending", zap.Error(err))
	}
}

func (s *Service) ScheduleMsgDeletionAtTime(adChanID string, channelID int64, messageID int, at time.Time) error {
	_, err := s.cron.AddFunc(
		fmt.Sprintf("%d %d %d %d *", at.Minute()+1, at.Hour(), at.Day(), int(at.Month())),
		func() {
			s.DeleteAdvertisement(adChanID, channelID, messageID)
		})

	return err
}
