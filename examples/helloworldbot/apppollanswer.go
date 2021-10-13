package main

import (
	"encoding/json"
	"fmt"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) PollAnswerHandler(update *go_telegram_bot_api.Update) {
	j, _ := json.Marshal(update)
	fmt.Println("PollAnswer update", string(j))
}
