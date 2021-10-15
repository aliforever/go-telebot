package main

import (
	"fmt"

	"github.com/aliforever/go-telebot/examples/helloworldbot/app"

	"github.com/aliforever/go-telebot"
)

func main() {
	token := "796493295:AAE3EGLAnba_XAsp_ts3sbPTHpW3nitBc4s"

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
