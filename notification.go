package dockerdiskwatcher

import (
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var BotToken string
var AdminId string
var AdminIdInt int64
var lastMessage time.Time
var Bot *tgbotapi.BotAPI

func SendTelegramMessage(message string) {
	AdminId = os.Getenv("TELEGRAM_ADMINID")
	BotToken = os.Getenv("TELEGRAM_APITOKEN")
	if AdminId == "" {
		panic("TELEGRAM_ADMINID is not set")
	}
	if BotToken == "" {
		panic("TELEGRAM_APITOKEN is not set")
	}
	var err error
	Bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}
	if time.Since(lastMessage).Seconds() < 60 {
		return
	}
	Bot.Send(tgbotapi.NewMessage(AdminIdInt, message))
	lastMessage = time.Now()
}
