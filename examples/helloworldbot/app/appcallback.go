package app

import (
	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) CallbackQueryHandler(update *go_telegram_bot_api.Update) {
	data := update.CallbackQuery.Data
	app.Send(app.AnswerCallbackQuery().SetCallbackQueryId(update.CallbackQuery.Id).SetText(data))
}
