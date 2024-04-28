package transport

import (
	"advertiser/shared/pkg/service/repo/models"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMe(t *testing.T) {
	origin := "фффlo, I am the link! nice to meet you"
	r := []rune(origin)
	res := r[16:20]
	fmt.Println(string(res))
}

func TestComposeAdMessage(t *testing.T) {
	res := composeText(models.Advertisement{
		MsgText: "I want to test посилання українською",
		MsgEntities: []models.MsgEntity{
			{
				Type:   "text_link",
				Offset: 25,
				Length: 11,
				URL:    "https://www.google.com",
			},
		},
	})

	assert.Equal(t, "I want to test посилання <a href=\"https://www.google.com\">українською</a>", res)
}
