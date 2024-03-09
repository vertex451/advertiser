package transport

import (
	"advertiser/shared/pkg/service/transport"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func (s *Service) editTopicsPrompt(respondTo int64, rawChannelID string) *transport.Msg {
	channelID, err := strconv.ParseInt(rawChannelID, 10, 64)
	if err != nil {
		zap.L().Panic("failed to parse string to int64")
	}

	msg := tgbotapi.NewMessage(respondTo, fmt.Sprintf(`
Choose topics from the list:
%s
`, strings.Join(s.uc.AllTopics(), ", ")))

	s.setState(respondTo, stateData{
		state:     StateEditTopics,
		channelID: channelID,
	})

	return &transport.Msg{
		Msg: msg,
	}
}
