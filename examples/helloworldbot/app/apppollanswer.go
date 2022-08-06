package app

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/aliforever/go-telegram-bot-api"
)

func (app App) PollAnswerHandler(update *tgbotapi.Update) {
	j, _ := json.Marshal(update)
	fmt.Println("PollAnswer update", string(j))
}
