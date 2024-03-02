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
	uc.repo.CreateAdvertisementChannelEntries(adsReadyToPost)
}

func (uc *UseCase) GetAdsChannelByStatus(status models.AdChanStatus) (map[int64][]models.AdvertisementChannel, error) {
	res, err := uc.repo.GetAdsChannelByStatus(status)
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

func (uc *UseCase) UpdateAdChanStatus(adChannelID string, status models.AdChanStatus) error {
	return uc.repo.UpdateAdChanStatus(adChannelID, status)
}

func (uc *UseCase) SetAdChanMessageID(adChanID string, msgID int) error {
	return uc.repo.SetAdChanMessageID(adChanID, msgID)
}
