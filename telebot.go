package telebot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliforever/go-telegram-bot-api"
	"reflect"
	"runtime/debug"

	log "github.com/sirupsen/logrus"

	"github.com/aliforever/go-telegram-bot-api/structs"
)

type Bot struct {
	token          string
	botInterface   interface{}
	api            *tgbotapi.TelegramBot
	stateStorage   UserStateStorage
	reflectType    reflect.Type
	updateHandlers *updateHandlers
	rateLimiter    RateLimiter
	isOfAPIType    bool
}

func NewBot(token string, app interface{}, options *BotOptions) (bot *Bot, api *tgbotapi.TelegramBot, err error) {
	if reflect.ValueOf(app).Kind() == reflect.Ptr {
		appName := reflect.ValueOf(app).Elem().Type().Name()
		err = errors.New(fmt.Sprintf("pass_app_without_pointer_as_%s{}_not_&%s{}", appName, appName))
		return
	}
	api, err = tgbotapi.NewTelegramBot(token)
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

func (bot *Bot) SetRateLimiter(limiter RateLimiter) {
	bot.rateLimiter = limiter // commit
}

func (bot *Bot) updateReplyStateNotExists(update tgbotapi.Update, state string) {
	message := update.Message
	if message == nil {
		message = update.EditedMessage
	}
	if message != nil {
		bot.stateStorage.SetUserState(message.Chat.Id, "Welcome")
	}
	j, _ := json.Marshal(update)
	log.Errorf("Handler for State: %s was not found!\n%s", state, string(j))
	return
}

func (bot *Bot) updateReplyStateInternalError(update tgbotapi.Update) {
	j, _ := json.Marshal(update)
	log.Errorf("Error getting user state: %s. For update: %s", update, string(j))
	return
}

func (bot *Bot) newApp() reflect.Value {
	app := reflect.New(bot.reflectType)

	return app
}

func (bot *Bot) appWithUpdate(app reflect.Value, update *tgbotapi.Update, defaultRecipient *int64) []reflect.Value {
	api := *bot.api

	if bot.isOfAPIType {
		if defaultRecipient != nil {
			api.SetRecipientChatId(*defaultRecipient)
		}
		app.Elem().Field(0).Set(reflect.ValueOf(api))
	}
	return []reflect.Value{app.Elem(), reflect.ValueOf(update)}
}

func (bot *Bot) newAppWithUpdate(defaultRecipientId *int64, update *tgbotapi.Update) []reflect.Value {
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

func (bot *Bot) invoke(app reflect.Value, update tgbotapi.Update, method string, isSwitched bool) {
	if app.MethodByName(method).Kind() == reflect.Invalid {
		bot.updateReplyStateNotExists(update, method)
		return
	}
	values := app.MethodByName(method).Call([]reflect.Value{reflect.ValueOf(&update), reflect.ValueOf(isSwitched)})
	if len(values) == 1 {
		if val, ok := values[0].Interface().(string); ok {
			if val != "" {
				bot.stateStorage.SetUserState(update.Message.Chat.Id, val)
				bot.invoke(app, update, val, true)
			}
		}
	}
}

func (bot *Bot) invokeMiddleware(app reflect.Value, update *tgbotapi.Update) (ignoreUpdate bool) {
	middlewareMethod := app.MethodByName("Middleware")
	if middlewareMethod.Kind() == reflect.Invalid {
		return
	}
	values := middlewareMethod.Call([]reflect.Value{reflect.ValueOf(update)})
	if len(values) == 1 {
		ignoreUpdate, _ = values[0].Interface().(bool)
	}
	return
}

func (bot *Bot) processUpdate(update tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			// log.Errorf("recovered from panic: %s\n%s", err, debug.Stack())
			fmt.Println(err)
			fmt.Println(string(debug.Stack()))
		}
	}()

	app := bot.newApp()
	api := *bot.api

	if update.From() != nil {
		api.SetRecipientChatId(update.From().Id)
	}

	app.Elem().Field(0).Set(reflect.ValueOf(api))

	if ignoreUpdate := bot.invokeMiddleware(app, &update); ignoreUpdate {
		return
	}

	var message *structs.Message
	if update.Message != nil {
		message = update.Message
	} else if update.EditedMessage != nil {
		message = update.EditedMessage
	}
	if message != nil && message.Chat.Type == "private" {
		state, err := bot.stateStorage.UserState(message.Chat.Id)
		if err != nil {
			bot.updateReplyStateInternalError(update)
			return
		}
		// _, exists := bot.reflectType.MethodByName(state)
		if app.MethodByName(state).Kind() == reflect.Invalid {
			bot.updateReplyStateNotExists(update, state)
			return
		}
		/*api := *bot.api

		 */
		apiField := app.Elem().Field(0)
		api.SetRecipientChatId(message.Chat.Id)
		apiField.Set(reflect.ValueOf(api))
		bot.invoke(app, update, state, false)
		return
	}
	if bot.updateHandlers == nil {
		return
	}
	bot.updateHandlers.processUpdate(app, update)
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
		if bot.rateLimiter != nil && bot.rateLimiter.ShouldLimitUpdate(&update) {
			continue
		}
		go bot.processUpdate(update)
	}
	return
}
