package main

import (
	"fmt"

	"github.com/aliforever/go-telebot"
)

func main() {
	token := ""

	var (
		bot *telebot.Bot
		err error
	)

	bot, _, err = telebot.NewBot(token, App{}, nil)
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
