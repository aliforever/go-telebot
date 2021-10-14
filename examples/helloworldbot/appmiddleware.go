package main

import (
	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) Middleware(update *go_telegram_bot_api.Update) (ignoreUpdate bool) {
	if update.Message != nil && update.Message.Text == "ignore" {
		ignoreUpdate = true
	}
	return
}
