package transport

import (
	"advertiser/channel_owner/internal/service/listener"
	"advertiser/channel_owner/internal/service/listener/transport"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

const (
	EveryTenSeconds = "*/10 * * * *"
)

type WriterService struct {
	uc       listener.UseCase
	tgBotApi *tgbotapi.BotAPI
}

func New(uc listener.UseCase, tgBotApi *tgbotapi.BotAPI) *WriterService {
	ws := &WriterService{
		tgBotApi: tgBotApi,
		uc:       uc,
	}

	go func() {
		for {
			ws.CheckForNewAds()
			time.Sleep(10 * time.Second)
		}
	}()

	return ws
}

// CheckForNewAds
// 1. query db, find new ads
// 2. send ads notification to channel admins
// 3. collect feedback and prepare for posting ad. (via cronJob?)
// 4. Post ad.
// 5. On expiration replace ad with thumbnail.
// 6. Collect statistics?
func (ws *WriterService) CheckForNewAds() {

	fmt.Println("###  CheckForNewAds()")
	res, err := ws.uc.CheckForNewAds()
	if err != nil {
		zap.L().Error("failed to get moderation receivers", zap.Error(err))
	}

	var text string
	var msg tgbotapi.MessageConfig
	for _, receiver := range res {
		text = fmt.Sprintf(`
Approve posting the following advertisement in your channel %s:
%s
`, receiver.Title, receiver.AdMessage)
		msg = tgbotapi.NewMessage(receiver.UserID, text)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("Approve"), fmt.Sprintf("%s/%s", transport.ApproveAd, receiver.AdID)),
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("Reject"), fmt.Sprintf("%s/%s", transport.RejectAd, receiver.AdID)),
			))
		_, err = ws.tgBotApi.Send(msg)
		if err != nil {
			zap.L().Error("failed to send message to moderation", zap.Error(err))
		}
	}
}

//func (ws *WriterService) StartCronJob() {
//	fmt.Println("### 1")
//	c := cron.New(
//		cron.WithParser(
//			cron.NewParser(
//				cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)))
//
//	_, err := c.AddFunc("*/1 * * * *", func() {
//		ws.CheckForNewAds()
//	})
//	fmt.Println("### 2")
//	if err != nil {
//		log.Fatal("Error adding cron job:", err)
//		return
//	}
//
//	fmt.Println("### 3")
//
//	c.start()
//	fmt.Println("### 4")
//}
