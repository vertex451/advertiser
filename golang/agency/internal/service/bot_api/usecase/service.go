package usecase

import (
	"go.uber.org/zap"
	"tg-bot/internal/service/bot_api"
)

type UseCase struct {
	repo bot_api.Repo
	cache
}

type cache struct {
	topics map[string]int
}

func New(repo bot_api.Repo) *UseCase {
	topics, err := repo.AllTopicsWithCoverage()
	if err != nil {
		zap.L().Panic("failed to init UseCase", zap.Error(err))
	}

	return &UseCase{
		repo: repo,
		cache: cache{
			topics: topics,
		},
	}
}

func (uc *UseCase) updateTopicCache() error {
	topics, err := uc.repo.AllTopicsWithCoverage()
	if err != nil {
		return err
	}

	uc.cache.topics = topics

	return nil
}
