package usecase

import (
	"advertiser/shared/pkg/service/repo/models"
	"go.uber.org/zap"
)

func (uc *UseCase) CheckForNewAds() {
	adsReadyToPost, err := uc.repo.GetAdsOnModeration()
	if err != nil {
		zap.L().Error("failed to fetch new advertisement for further moderation",
			zap.Error(err))
	}

	// create new entries in advertisement_channel table
	uc.repo.CreateAdChanEntries(adsReadyToPost)
}

func (uc *UseCase) GetAdsChannelByStatus(status models.AdChanStatus) (map[int64][]models.AdvertisementChannel, error) {
	res, err := uc.repo.GetAdChannelByStatus(status)
	if err != nil {
		zap.L().Error("failed to get advertisement channel by status", zap.Error(err))
		return nil, err
	}

	m := make(map[int64][]models.AdvertisementChannel)
	for _, entry := range res {
		m[entry.ChannelOwnerID] = append(m[entry.ChannelOwnerID], entry)
	}

	return m, nil
}

func (uc *UseCase) UpdateAdChanEntry(channel models.AdvertisementChannel) error {
	return uc.repo.UpdateAdChanEntry(channel)
}

func (uc *UseCase) GetRunningAdvertisements() ([]*models.Advertisement, error) {
	return uc.repo.GetRunningAds()
}

func (uc *UseCase) UpdateAd(ad models.Advertisement) error {
	return uc.repo.UpdateAd(ad)
}
