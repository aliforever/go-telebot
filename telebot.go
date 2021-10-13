package telebot

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"

	log "github.com/sirupsen/logrus"

	"github.com/GoLibs/telegram-bot-api/structs"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

type Bot struct {
	token          string
	botInterface   interface{}
	api            *go_telegram_bot_api.TelegramBot
	stateStorage   UserStateStorage
	reflectType    reflect.Type
	updateHandlers *updateHandlers
	isOfAPIType    bool
}

func NewBot(token string, app interface{}, options *BotOptions) (bot *Bot, api *go_telegram_bot_api.TelegramBot, err error) {
	if reflect.ValueOf(app).Kind() == reflect.Ptr {
		appName := reflect.ValueOf(app).Elem().Type().Name()
		err = errors.New(fmt.Sprintf("pass_app_without_pointer_as_%s{}_not_&%s{}", appName, appName))
		return
	}
	api, err = go_telegram_bot_api.NewTelegramBot(token)
	if err != nil {
		return
	}

	t := reflect.TypeOf(app)

	defaultUserStateStorage := newStateStorage()
	bot = &Bot{
		token:        token,
		botInterface: app,
		api:          api,
		stateStorage: defaultUserStateStorage,
		reflectType:  t,
	}
	bot.updateHandlers = updateHandlersFromType(bot, t)
	var v = reflect.ValueOf(app)
	if reflect.ValueOf(app).Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.NumField() > 0 && v.Field(0).Type().AssignableTo(reflect.TypeOf(api).Elem()) {
		bot.isOfAPIType = true
	}
	return
}

func (bot *Bot) SetUserStateStorage(storage UserStateStorage) {
	bot.stateStorage = storage
}

func (bot *Bot) updateReplyStateNotExists(update *go_telegram_bot_api.Update, state string) {
	j, _ := json.Marshal(update)
	log.Errorf("Handler for State: %s was not found!\n%s", state, string(j))
	return
}

func (bot *Bot) updateReplyStateInternalError(update *go_telegram_bot_api.Update) {
	j, _ := json.Marshal(update)
	log.Errorf("Error getting user state: %s. For update: %s", update, string(j))
	return
}

func (bot *Bot) newAppWithUpdate(defaultRecipientId *int64, update *go_telegram_bot_api.Update) []reflect.Value {
	app := reflect.New(bot.reflectType)
	api := *bot.api
	if defaultRecipientId != nil {
		api.SetRecipientChatId(*defaultRecipientId)
	}
	if bot.isOfAPIType {
		app.Elem().Field(0).Set(reflect.ValueOf(api))
	}
	return []reflect.Value{app.Elem(), reflect.ValueOf(update)}
}

func (bot *Bot) invoke(update *go_telegram_bot_api.Update, method string, isSwitched bool) {
	app := reflect.New(bot.reflectType)
	api := *bot.api
	api.SetRecipientChatId(update.Message.Chat.Id)
	if bot.isOfAPIType {
		app.Elem().Field(0).Set(reflect.ValueOf(api))
	}
	fn := app.MethodByName(method)
	values := fn.Call([]reflect.Value{reflect.ValueOf(update), reflect.ValueOf(isSwitched)})
	if len(values) == 1 {
		/*if values[0].Type().Implements(reflect.TypeOf((*go_telegram_bot_api.Config)(nil)).Elem()) {
			if values[0].Interface() != nil {
				api.Send(values[0].Interface().(go_telegram_bot_api.Config))
			}
		}*/
		if val, ok := values[0].Interface().(string); ok {
			if val != "" {
				bot.stateStorage.SetUserState(update.Message.Chat.Id, val)
				bot.invoke(update, val, true)
			}
		}
	}
}

func (bot *Bot) processUpdate(update *go_telegram_bot_api.Update) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("recovered from panic: %s\n%s", err, debug.Stack())
		}
	}()
	var message *structs.Message
	if update.Message != nil {
		message = update.Message
	} else if update.EditedMessage != nil {
		message = update.EditedMessage
	}
	if message != nil && message.Chat.Type == "private" {
		state, err := bot.stateStorage.UserState(update.Message.Chat.Id)
		if err != nil {
			bot.updateReplyStateInternalError(update)
			return
		}
		_, exists := bot.reflectType.MethodByName(state)
		if !exists {
			bot.updateReplyStateNotExists(update, state)
			return
		}
		bot.invoke(update, state, false)
		return
	}
	if bot.updateHandlers == nil {
		return
	}
	bot.updateHandlers.processUpdate(update)
	return
}

func (bot *Bot) Poll() (err error) {
	var allowedUpdates []string = []string{"message", "edited_message"}
	if bot.updateHandlers != nil {
		allowedUpdates = bot.updateHandlers.allowedUpdates()
	}

	for update := range bot.api.GetUpdates().SetAllowedUpdates(allowedUpdates).LongPoll() {
		if err = update.Error(); err != nil {
			return
		}
		go bot.processUpdate(update)
	}
	return
}
