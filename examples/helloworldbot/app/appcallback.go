package app

func (app App) CallbackQueryHandler(update *tgbotapi.Update) {
	data := update.CallbackQuery.Data
	app.Send(app.AnswerCallbackQuery().SetCallbackQueryId(update.CallbackQuery.Id).SetText(data))
}
