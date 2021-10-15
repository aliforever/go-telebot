package telebot

import (
	"flag"
	"os"
	"strings"

	"github.com/aliforever/go-telebot/generator"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

var (
	newFlag  string
	langsStr string
)

func init() {
	flag.StringVar(&newFlag, "new", "", "--new=BOT_TOKEN_HERE")
	flag.StringVar(&langsStr, "langs", "", "--langs=Persian,Arabic,French")
	flag.Parse()

	var langsSlice []string
	if langsStr != "" {
		langsSlice = strings.Split(langsStr, ",")
	}
	if newFlag != "" {
		defer os.Exit(1)
		_, err := go_telegram_bot_api.NewTelegramBot(newFlag)
		if err != nil {
			log.Error("error creating new bot", err)
			return
		}
		err = generator.NewBot(newFlag, langsSlice)
		if err != nil {
			log.Error("error creating new bot", err)
			return
		}
	} else if langsStr != "" {
		// TODO: Create a new language file
	}
}
