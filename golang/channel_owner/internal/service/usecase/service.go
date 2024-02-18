package usecase

import (
	"advertiser/channel_owner/internal/service"
	"advertiser/channel_owner/internal/service/transport"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
)

type UseCase struct {
	repo service.Repo
	cache
}

type cache struct {
	topics map[string]struct{}
}

func New(repo service.Repo) *UseCase {
	topics, err := repo.AllTopics()
	if err != nil {
		zap.L().Panic("failed to init UseCase", zap.Error(err))
	}

	topicMap := make(map[string]struct{})
	for _, topic := range topics {
		topicMap[topic] = struct{}{}
	}

	return &UseCase{
		repo: repo,
		cache: cache{
			topics: topicMap,
		},
	}
}

func (uc *UseCase) updateTopicCache() error {
	topics, err := uc.repo.AllTopics()
	if err != nil {
		return err
	}

	topicMap := make(map[string]struct{})
	for _, topic := range topics {
		topicMap[topic] = struct{}{}
	}

	uc.cache.topics = topicMap

	return nil
}

func (uc *UseCase) validateTopics(input []string) error {
	var notRecognisedTopics []string
	var ok bool
	for _, inputTopic := range input {
		if _, ok = uc.topics[inputTopic]; !ok {
			notRecognisedTopics = append(notRecognisedTopics, inputTopic)
		}
	}

	if len(notRecognisedTopics) > 0 {
		return errors.Errorf("invalid topics: %s. Please use /%s command to list allowed topics",
			strings.Join(notRecognisedTopics, ", "), transport.AllTopics)
	}

	return nil
}
