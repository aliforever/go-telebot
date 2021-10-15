package app

import (
	"fmt"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) MessageTypeGroupHandler(update *go_telegram_bot_api.Update) {
	app.Send(app.Message().SetText(fmt.Sprintf("This message is from Group: %s", update.Message.Chat.Title)))
}
