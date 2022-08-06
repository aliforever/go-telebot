package app

import tgbotapi "github.com/aliforever/go-telegram-bot-api"

func (app App) CallbackQueryHandler(update *tgbotapi.Update) {
	data := update.CallbackQuery.Data
	app.Send(app.AnswerCallbackQuery().SetCallbackQueryId(update.CallbackQuery.Id).SetText(data))
}
