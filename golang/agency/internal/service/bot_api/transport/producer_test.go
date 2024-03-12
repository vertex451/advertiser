package transport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAndValidateCreateAdInput(t *testing.T) {
	msg := `
To create an advertisement, send a message in the following format:
Name: Stock market
TargetTopics: topic1, topic2, topic3
BudgetUSD: 100
CostPerView: 0.1
Message: Follow this [link](https://www.investing.com/) to find more about investments!
`
	ad, err := parseAndValidateCreateAdInput(
		"57be371d-6674-4c65-af14-4ce273917e25",
		"57be371d-6674-4c65-af14-4ce273917e25",
		msg)
	assert.Nil(t, err)
	assert.Equal(t, "Follow this [link](https://www.investing.com/) to find more about investments!", ad.Message)
}
