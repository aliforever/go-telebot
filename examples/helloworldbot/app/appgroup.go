package app

import (
	"fmt"
)

func (app App) MessageTypeGroupHandler(update *tgbotapi.Update) {
	app.Send(app.Message().SetText(fmt.Sprintf("This message is from Group: %s", update.Message.Chat.Title)))
}
