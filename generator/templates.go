package generator

import "fmt"

func mainTemplate(token string) (str string) {
	str = `token := "%s"

	var (
		bot *telebot.Bot
		err error
	)

	bot, _, err = telebot.NewBot(token, app.App{}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = bot.Poll()
	if err != nil {
		fmt.Println(err)
		return
	}`
	return fmt.Sprintf(str, token)
}

func appTemplate() (str string) {
	str = `package app

import (
	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

type App struct {
	go_telegram_bot_api.TelegramBot
}
`
	return
}

func appCallbackTemplate() (str string) {
	str = `package app

import (
	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) CallbackQueryHandler(update *go_telegram_bot_api.Update) {
	data := update.CallbackQuery.Data
	app.Send(app.AnswerCallbackQuery().SetCallbackQueryId(update.CallbackQuery.Id).SetText(data))
}`
	return
}

func appChatMemberTemplate() (str string) {
	str = `package app

import (
	"encoding/json"
	"fmt"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) ChatMemberHandler(update *go_telegram_bot_api.Update) {
	j, _ := json.Marshal(update)
	fmt.Println("chat_member update", string(j))
}
`
	return
}

func appPollAnswer() (str string) {
	str = `package app

import (
	"encoding/json"
	"fmt"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) PollAnswerHandler(update *go_telegram_bot_api.Update) {
	j, _ := json.Marshal(update)
	fmt.Println("PollAnswer update", string(j))
}
`
	return
}

func appGroupTemplate() (str string) {
	str = `package app

import (
	"fmt"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) MessageTypeGroupHandler(update *go_telegram_bot_api.Update) {
	app.Send(app.Message().SetText(fmt.Sprintf("This message is from Group: %s", update.Message.Chat.Title)))
}
`
	return
}

func appMyChatMemberTemplate() (str string) {
	str = `package app

import (
	"encoding/json"
	"fmt"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) MyChatMemberHandler(update *go_telegram_bot_api.Update) {
	j, _ := json.Marshal(update)
	fmt.Println("chat_member update", string(j))
}
`
	return
}

func appChannelTemplate() (str string) {
	str = `package app

import (
	"fmt"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) ChannelPostHandler(update *go_telegram_bot_api.Update) {
	app.Send(app.Message().SetText(fmt.Sprintf("This message is from Channel: %s", update.ChannelPost.Chat.Title)))
}
`
	return
}

func appPrivateTemplate() (str string) {
	str = `package app

import (
	"fmt"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

func (app App) Welcome(update *go_telegram_bot_api.Update, isSwitched bool) (newState string) {
	if update.Message != nil && update.Message.Text != "" {
		if !isSwitched {
			if update.Message.Text == "Hello" {
				newState = "Hello"
				return
			}
		}
	}

	keyboard := app.Tools.Keyboards.NewReplyKeyboardFromSlicesOfStrings([][]string{
		{"Hello"},
	}).SetResizeKeyboard(true)
	app.Send(app.Message().SetText(fmt.Sprintf("Hello World!\nYour name is: %s", update.Message.Chat.FirstName)).SetReplyMarkup(keyboard))
	return
}

func (app App) Hello(update *go_telegram_bot_api.Update, isSwitched bool) (newState string) {
	if !isSwitched {
		if update.Message.Text == "Back" {
			newState = "Welcome"
			return
		}
	}
	keyboard := app.Tools.Keyboards.NewReplyKeyboardFromSlicesOfStrings([][]string{
		{"Back"},
	}).SetResizeKeyboard(true)
	app.Send(app.Message().SetText(fmt.Sprintf("You are in Hello Menu %s", update.Message.Chat.FirstName)).SetReplyMarkup(keyboard))
	return
}
`
	return
}
