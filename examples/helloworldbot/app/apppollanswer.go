package app

import (
	"encoding/json"
	"fmt"
)

func (app App) PollAnswerHandler(update *tgbotapi.Update) {
	j, _ := json.Marshal(update)
	fmt.Println("PollAnswer update", string(j))
}
