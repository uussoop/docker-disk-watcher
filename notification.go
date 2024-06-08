package dockerdiskwatcher

import (
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var BotToken string
var AdminId string
var AdminIdInt int64
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
	AdminIdInt, err := strconv.ParseInt(AdminId, 10, 64)
	if err != nil {
		panic(err)
	}
	Bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}
	Bot.Send(tgbotapi.NewMessage(AdminIdInt, message))

}
