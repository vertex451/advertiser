package usecase

import (
	"advertiser/shared/pkg/service/constants"
	"github.com/pkg/errors"
	"strings"
)

func ValidateTopics(supportedTopics map[string]int, input []string) error {
	var notRecognisedTopics []string
	var ok bool
	for _, inputTopic := range input {
		if _, ok = supportedTopics[inputTopic]; !ok {
			notRecognisedTopics = append(notRecognisedTopics, inputTopic)
		}
	}

	if len(notRecognisedTopics) > 0 {
		return errors.Errorf("invalid topics: %s. Please use /%s command to list allowed topics",
			strings.Join(notRecognisedTopics, ", "), constants.AllTopics)
	}

	return nil
}
