package app

import (
	"fmt"
	tgbotapi "github.com/aliforever/go-telegram-bot-api"
)

func (app *App) Middleware(update *tgbotapi.Update) (ignoreUpdate bool) {
	fmt.Println("middleware is called")
	app.CustomField = 2
	if update.Message != nil && update.Message.Text == "ignore" {
		ignoreUpdate = true
	}
	return
}
