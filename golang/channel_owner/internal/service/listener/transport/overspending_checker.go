package transport

import (
	"advertiser/shared/pkg/service/repo/models"
	"go.uber.org/zap"
)

const (
	DeletePostThreshold = 0.8
)

func (s *Service) PreventOverspending() {
	var err error
	runningAds, err := s.uc.GetRunningAdvertisements()
	if err != nil {
		zap.L().Error("failed to get running advertisements", zap.Error(err))
		return
	}

	for _, ad := range runningAds {
		s.CalculateTotalCostForSingleAd(*ad)
	}
}

func (s *Service) CalculateTotalCostForSingleAd(ad models.Advertisement) {
	var err error
	var totalAdViews, msgViews int
	for _, adChan := range ad.AdsChannel {
		if adChan.Status == models.AdChanPosted {
			msgViews, err = s.tgApi.GetMessageViews(adChan.ChannelHandle, adChan.MessageID)
			if err != nil {
				zap.L().Error("failed to get message views", zap.Error(err))
				// TODO add backoff
				continue
			}
			totalAdViews += msgViews
		}
	}

	deleteAdErr := false
	if float32(totalAdViews)*ad.CostPerView >= float32(ad.Budget)*DeletePostThreshold {
		for _, adChan := range ad.AdsChannel {
			err = s.DeleteAdvertisement(adChan.ID.String(), adChan.ChannelID, adChan.MessageID)
			if err != nil {
				deleteAdErr = true
			}
		}
		if !deleteAdErr {
			ad.Status = models.AdsStatusOutOfBudget
		}
	}

	ad.TotalViews = totalAdViews
	ad.TotalCost = float32(totalAdViews) * ad.CostPerView

	err = s.uc.UpdateAd(ad)
}
