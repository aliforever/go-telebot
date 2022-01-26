package app

import (
	"encoding/json"
	"fmt"
)

func (app App) ChatMemberHandler(update *tgbotapi.Update) {
	j, _ := json.Marshal(update)
	fmt.Println("chat_member update", string(j))
}
