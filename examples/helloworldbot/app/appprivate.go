package app

import (
	"encoding/json"
	"fmt"

	"github.com/aliforever/go-telegram-bot-api/structs"
)

func (app *App) Welcome(update *tgbotapi.Update, isSwitched bool) (newState string) {
	var message *structs.Message
	if update.Message != nil {
		message = update.Message
	} else if update.EditedMessage != nil {
		message = update.EditedMessage
	}
	if message == nil {
		j, _ := json.Marshal(update)
		fmt.Println("unknown", string(j))
		return
	}
	app.Send(app.Message().SetText("Hi").SetReplyToMessageId(message.MessageId))
	if !isSwitched {
		if message.Text == "Hello" {
			app.Send(app.Message().SetText("Hi Bruh!"))
			return
		}
		if message.Text == "Bye" {
			app.Send(app.Message().SetText("Bye Bruh!"))
			newState = "Bye"
			return
		}
		if message.Text == "Inline" {
			keyboard := app.Tools.Keyboards.NewInlineKeyboardFromSlicesOfMaps([][]map[string]string{
				{
					{"text": "Click on Me!", "callback_data": "clicked_here"},
				},
			})
			app.Send(app.Message().SetText("Click!").SetReplyMarkup(keyboard))
			return
		}
		if message.Text == "Poll" {
			app.Send(app.Poll().SetQuestion("How Are You Today?").SetOptions([]string{
				"Cool",
				"Fine",
				"Well",
			}).DisableAnonymous())
			return
		}
	}
	keyboard := app.Tools.Keyboards.NewReplyKeyboardFromSlicesOfStrings([][]string{
		{"Hello", "Bye"},
		{"Inline", "Poll"},
	}).SetResizeKeyboard(true)
	app.Send(app.Message().SetText(fmt.Sprintf("Hello World!\nYour name is: %s", update.Message.Chat.FirstName)).SetReplyMarkup(keyboard))
	return
}

func (app App) Bye(update *tgbotapi.Update, isSwitched bool) (newState string) {
	if !isSwitched {
		if update.Message.Text == "Back" {
			fmt.Println(app.CustomField)
			newState = "Welcome"
			return
		}
	}
	keyboard := app.Tools.Keyboards.NewReplyKeyboardFromSlicesOfStrings([][]string{
		{"Back"},
	}).SetResizeKeyboard(true)
	app.Send(app.Message().SetText(fmt.Sprintf("You are in Bye Menu %s", update.Message.Chat.FirstName)).SetReplyMarkup(keyboard))
	return
}
