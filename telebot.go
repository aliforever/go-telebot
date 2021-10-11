package telebot

import (
	"reflect"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

type Bot struct {
	token        string
	botInterface interface{}
	api          *go_telegram_bot_api.TelegramBot
	stateStorage UserStateStorage
	reflectType  reflect.Type
}

func NewBot(token string, botInterface interface{}, options *BotOptions) (bot *Bot, api *go_telegram_bot_api.TelegramBot, err error) {
	api, err = go_telegram_bot_api.NewTelegramBot(token)
	if err != nil {
		return
	}

	t := reflect.TypeOf(botInterface)
	defaultUserStateStorage := newStateStorage()
	bot = &Bot{
		token:        token,
		botInterface: botInterface,
		api:          api,
		stateStorage: defaultUserStateStorage,
		reflectType:  t,
	}
	return
}

func (bot *Bot) SetUserStateStorage(storage UserStateStorage) {
	bot.stateStorage = storage
}

func (bot *Bot) updateReplyStateNotExists(update *go_telegram_bot_api.Update) {
	// TODO: Make it so they use their own message for this error
	bot.api.Send(bot.api.Message().SetChatId(update.Message.Chat.Id).SetText("Oops, you're lost in the bot!"))
	return
}

func (bot *Bot) updateReplyStateInternalError(update *go_telegram_bot_api.Update) {
	// TODO: Make it so they use their own message for this error
	bot.api.Send(bot.api.Message().SetChatId(update.Message.Chat.Id).SetText("Oops, Internal Error!"))
	return
}

func (bot *Bot) invoke(update *go_telegram_bot_api.Update, method string, isSwitched bool) {
	app := reflect.New(bot.reflectType)
	api := *bot.api
	api.SetRecipientChatId(update.Message.Chat.Id)
	if app.Elem().NumField() > 0 && app.Elem().Field(0).Type().AssignableTo(reflect.TypeOf(api)) {
		app.Elem().Field(0).Set(reflect.ValueOf(api))
	}
	fn := app.MethodByName(method)
	values := fn.Call([]reflect.Value{reflect.ValueOf(update), reflect.ValueOf(isSwitched)})
	if len(values) == 2 {
		if values[0].Type().Implements(reflect.TypeOf((*go_telegram_bot_api.Config)(nil)).Elem()) {
			if values[0].Interface() != nil {
				api.Send(values[0].Interface().(go_telegram_bot_api.Config))
			}
		}
		if val, ok := values[1].Interface().(string); ok {
			if val != "" {
				bot.stateStorage.SetUserState(update.Message.Chat.Id, val)
				bot.invoke(update, val, true)
			}
		}
	}
}

func (bot *Bot) processUpdate(update *go_telegram_bot_api.Update) {
	if update.Message != nil && update.Message.Chat.Type == "private" {
		state, err := bot.stateStorage.UserState(update.Message.Chat.Id)
		if err != nil {
			bot.updateReplyStateInternalError(update)
			return
		}
		_, exists := bot.reflectType.MethodByName(state)
		if !exists {
			bot.updateReplyStateNotExists(update)
			return
		}
		bot.invoke(update, state, false)
	}
	return
}

func (bot *Bot) Poll() (err error) {
	for update := range bot.api.GetUpdates().LongPoll() {
		if err = update.Error(); err != nil {
			return
		}
		go bot.processUpdate(update)
	}
	return
}
