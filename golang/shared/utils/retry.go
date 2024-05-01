package utils

import "time"

const (
	DefaultMaxRetries      = 3           // Default maximum number of retries
	DefaultBaseBackoffTime = time.Second // Default base backoff duration
)

// RetryWithBackoff retries a function with an exponential backoff in case of errors.
func RetryWithBackoff(retryFunc func() (interface{}, error), maxRetries int, baseBackoffTime time.Duration) (interface{}, error) {
	var err error
	var result interface{}
	backoffDuration := baseBackoffTime

	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, err = retryFunc()
		if err == nil {
			return result, nil // Return the result if the function succeeds
		}

		// If this is the last attempt, return the error
		if attempt == maxRetries {
			break
		}

		// Exponential backoff
		time.Sleep(backoffDuration)
		backoffDuration *= 2
	}

	return nil, err // Return the final error after all attempts
}
