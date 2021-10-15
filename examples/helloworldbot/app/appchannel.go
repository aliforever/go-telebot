package app

import (
	"fmt"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) ChannelPostHandler(update *go_telegram_bot_api.Update) {
	app.Send(app.Message().SetText(fmt.Sprintf("This message is from Channel: %s", update.ChannelPost.Chat.Title)))
}
