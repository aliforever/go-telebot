package main

import (
	"fmt"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

type App struct {
	go_telegram_bot_api.TelegramBot
}

func (app App) Welcome(update *go_telegram_bot_api.Update, isSwitched bool) (cfg go_telegram_bot_api.Config, newState string) {
	if !isSwitched {
		if update.Message.Text == "Hello" {
			cfg = app.Message().SetText("Hi Bruh!")
			return
		}
		if update.Message.Text == "Bye" {
			cfg = app.Message().SetText("Bye Bruh!")
			newState = "Bye"
			return
		}
	}
	keyboard := app.Tools.Keyboards.NewReplyKeyboardFromSlicesOfStrings([][]string{
		{"Hello", "Bye"},
	}).SetResizeKeyboard(true)
	cfg = app.Message().SetText(fmt.Sprintf("Hello World!\nYour name is: %s", update.Message.Chat.FirstName)).SetReplyMarkup(keyboard)
	return
}

func (app App) Bye(update *go_telegram_bot_api.Update, isSwitched bool) (cfg go_telegram_bot_api.Config, newState string) {
	if !isSwitched {
		if update.Message.Text == "Back" {
			newState = "Welcome"
			return
		}
	}
	keyboard := app.Tools.Keyboards.NewReplyKeyboardFromSlicesOfStrings([][]string{
		{"Back"},
	}).SetResizeKeyboard(true)
	cfg = app.Message().SetText(fmt.Sprintf("You are in Bye Menu %s", update.Message.Chat.FirstName)).SetReplyMarkup(keyboard)
	return
}
