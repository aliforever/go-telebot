package generator

import (
	"fmt"
)

func langInterface() (str string) {
	str = `package langs

type Language interface {
	Flag() string
	ChooseLanguageText() string
	WelcomeMenu() string
}
`
	return
}

func langFile(name string) (str string) {
	str = `package langs

type %s struct {
}

func (English) Flag() string {
	return ""
}

func (English) WelcomeMenu() string {
	return ""
}

func (English) ChooseLanguageText() string {
	return ""
}
`
	return fmt.Sprintf(str, name)
}

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
	 "github.com/aliforever/go-telegram-bot-api"
)

type App struct {
	tgbotapi.TelegramBot
}
`
	return
}

func appCallbackTemplate() (str string) {
	str = `package app

import (
	 "github.com/aliforever/go-telegram-bot-api"
)

func (app App) CallbackQueryHandler(update *tgbotapi.Update) {
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

	 "github.com/aliforever/go-telegram-bot-api"
)

func (app App) ChatMemberHandler(update *tgbotapi.Update) {
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

	 "github.com/aliforever/go-telegram-bot-api"
)

func (app App) PollAnswerHandler(update *tgbotapi.Update) {
	j, _ := json.Marshal(update)
	fmt.Println("PollAnswer update", string(j))
}
`
	return
}

func appMiddlewareTemplate() (str string) {
	str = `package app

import (
	 "github.com/aliforever/go-telegram-bot-api"
)

func (app App) Middleware(update *tgbotapi.Update) (ignoreUpdate bool) {
	return
}
`
	return
}

func appGroupTemplate() (str string) {
	str = `package app

import (
	"fmt"

	 "github.com/aliforever/go-telegram-bot-api"
)

func (app App) MessageTypeGroupHandler(update *tgbotapi.Update) {
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

	 "github.com/aliforever/go-telegram-bot-api"
)

func (app App) MyChatMemberHandler(update *tgbotapi.Update) {
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

	 "github.com/aliforever/go-telegram-bot-api"
)

func (app App) ChannelPostHandler(update *tgbotapi.Update) {
	app.Send(app.Message().SetText(fmt.Sprintf("This message is from Channel: %s", update.ChannelPost.Chat.Title)))
}
`
	return
}

func appPrivateTemplate() (str string) {
	str = `package app

import (
	"fmt"

	 "github.com/aliforever/go-telegram-bot-api"
)

func (app App) Welcome(update *tgbotapi.Update, isSwitched bool) (newState string) {
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

func (app App) Hello(update *tgbotapi.Update, isSwitched bool) (newState string) {
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
