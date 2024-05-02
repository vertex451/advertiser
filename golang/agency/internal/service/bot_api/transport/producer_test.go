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
CostPerMile: 0.1
Message: Follow this [link](https://www.investing.com/) to find more about investments!
`
	ad, err := parseAndValidateCreateAdInput("57be371d-6674-4c65-af14-4ce273917e25", msg)
	assert.Nil(t, err)
	assert.Equal(t, "Follow this [link](https://www.investing.com/) to find more about investments!", ad.MsgText)
}

func TestComposeAdMessage(t *testing.T) {
	//	msg := transport.ComposeAdMessage(&models.Advertisement{
	//		MsgText: `Welcome to the advertisement bot!
	//
	//You can create rich messages using the following formatting options:
	//bold, italic, strikethrough, underline,  link, monospace or even spoiler!
	//You can attach the image as well!
	//
	//Just send me a message as you usually do and I will keep the formatting.`,
	//		MsgEntities: []models.MsgEntity{
	//			{Type: "bold", Offset: 0, Length: 33},
	//			{Type: "bold", Offset: 104, Length: 4},
	//			{Type: "italic", Offset: 110, Length: 6},
	//			{Type: "strikethrough", Offset: 118, Length: 13},
	//			{Type: "underline", Offset: 133, Length: 9},
	//			{Type: "text_link", Offset: 145, Length: 4, URL: "https://www.google.com/"},
	//			{Type: "code", Offset: 151, Length: 9},
	//			{Type: "spoiler", Offset: 169, Length: 7},
	//		},
	//	})
	//	assert.Equal(t, `<b>Welcome to the advertisement bot!</b>
	//
	//You can create rich messages using the following formatting options:
	//<b>bold</b>, <i>italic</i>, <strike>strikethrough</strike>, <u>underline</u>,  <a href="https://www.google.com/">link</a>, <code>monospace</code> or even <tg-spoiler>spoiler</tg-spoiler>!
	//You can attach the image as well!
	//
	//Just send me a message as you usually do and I will keep the formatting.`, msg)
}
