package usecase

import (
	"advertiser/channel_owner/internal/service/listener/repo/postgresql/types"
)

func (uc *UseCase) CheckForNewAds() ([]types.GetAdsOnModerationResult, error) {

	return uc.repo.GetAdsOnModeration()
}
