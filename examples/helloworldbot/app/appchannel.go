package app

import (
	"fmt"
	tgbotapi "github.com/aliforever/go-telegram-bot-api"
)

func (app App) ChannelPostHandler(update *tgbotapi.Update) {
	app.Send(app.Message().SetText(fmt.Sprintf("This message is from Channel: %s", update.ChannelPost.Chat.Title)))
}
