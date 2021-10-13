package main

import (
	"fmt"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) CallbackQueryHandler(update *go_telegram_bot_api.Update) {
	defer func() {
		app.Send(app.AnswerCallbackQuery().SetCallbackQueryId(update.CallbackQuery.Id))
	}()
	data := update.CallbackQuery.Data
	fmt.Println(data)
}
