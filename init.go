package telebot

import (
	"flag"
	tgbotapi "github.com/aliforever/go-telegram-bot-api"
	"os"
	"strings"

	"github.com/aliforever/go-telebot/generator"

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
		_, err := tgbotapi.New(newFlag)
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
