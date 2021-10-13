package telebot

import (
	"flag"
	"os"

	"github.com/aliforever/go-telebot/generator"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

var (
	newFlag string
)

func init() {
	flag.StringVar(&newFlag, "new", "", "--new=BOT_TOKEN_HERE")
	flag.Parse()

	if newFlag != "" {
		defer os.Exit(1)
		_, err := go_telegram_bot_api.NewTelegramBot(newFlag)
		if err != nil {
			log.Error("error creating new bot", err)
			return
		}
		err = generator.NewBot(newFlag)
		if err != nil {
			log.Error("error creating new bot", err)
			return
		}
	}
}
