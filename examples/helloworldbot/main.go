package main

import (
	"fmt"

	"github.com/aliforever/go-telebot/examples/helloworldbot/app"

	"github.com/aliforever/go-telebot"
)

func main() {
	token := "357820880:AAHYTDRZtSCgcVgkSAlhRuntE-8gJNnU2IE"

	var (
		bot *telebot.Bot
		err error
	)

	bot, _, err = telebot.NewBot(token, app.App{}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = bot.Poll()
	if err != nil {
		fmt.Println(err)
		return
	}
}
