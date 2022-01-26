package app

import (
	"fmt"
)

func (app App) ChannelPostHandler(update *tgbotapi.Update) {
	app.Send(app.Message().SetText(fmt.Sprintf("This message is from Channel: %s", update.ChannelPost.Chat.Title)))
}
