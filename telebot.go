package telebot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime/debug"

	"github.com/aliforever/go-telegram-bot-api"

	log "github.com/sirupsen/logrus"

	"github.com/aliforever/go-telegram-bot-api/structs"
)

type Bot struct {
	token          string
	botInterface   interface{}
	api            *tgbotapi.TelegramBot
	reflectType    reflect.Type
	updateHandlers *updateHandlers
	rateLimiter    RateLimiter
	isOfAPIType    bool
	options        *BotOptions
}

func NewBot(token string, app interface{}, options *BotOptions) (bot *Bot, api *tgbotapi.TelegramBot, err error) {
	if reflect.ValueOf(app).Kind() == reflect.Ptr {
		appName := reflect.ValueOf(app).Elem().Type().Name()
		err = errors.New(fmt.Sprintf("pass_app_without_pointer_as_%s{}_not_&%s{}", appName, appName))
		return
	}

	api, err = tgbotapi.New(token)
	if err != nil {
		return
	}

	t := reflect.TypeOf(app)

	if options == nil {
		options = NewOptions()
	}

	if options.stateStorage == nil {
		options.stateStorage = newStateStorage()
	}

	bot = &Bot{
		token:        token,
		botInterface: app,
		api:          api,
		reflectType:  t,
		options:      options,
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

func (bot *Bot) SetRateLimiter(limiter RateLimiter) {
	bot.rateLimiter = limiter // commit
}

func (bot *Bot) updateReplyStateNotExists(update tgbotapi.Update, state string) {
	message := update.Message
	if message == nil {
		message = update.EditedMessage
	}
	if message != nil {
		bot.options.stateStorage.SetUserState(message.Chat.Id, "Welcome")
	}
	j, _ := json.Marshal(update)
	log.Errorf("Handler for State: %s was not found!\n%s", state, string(j))
	return
}

func (bot *Bot) updateReplyStateInternalError(update tgbotapi.Update, err error) {
	if bot.options != nil && bot.options.logger != nil {
		bot.options.logger.Error(
			"Error getting user state",
			slog.Any("update", update),
			slog.String("err", err.Error()),
		)
	}

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
				bot.options.stateStorage.SetUserState(update.Message.Chat.Id, val)
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
			if bot.options != nil && bot.options.logger != nil {
				bot.options.logger.Error(
					"recovered from panic",
					slog.Any("update", update),
					slog.Any("err", err),
					slog.String("stack", string(debug.Stack())),
				)
			}
		}
	}()

	app := bot.newApp()
	api := *bot.api

	if update.From() != nil {
		api.SetRecipientChatId(update.From().Id)
		if bot.options.useStorage != nil {
			bot.options.useStorage.Store(update.From())
		}
	}

	app.Elem().Field(0).Set(reflect.ValueOf(api))

	fmt.Println("here 1", update.PreCheckoutQuery == nil)

	if ignoreUpdate := bot.invokeMiddleware(app, &update); ignoreUpdate {
		return
	}

	fmt.Println("here 2", update.PreCheckoutQuery == nil)

	var message *structs.Message
	if update.Message != nil {
		message = update.Message
	} else if update.EditedMessage != nil {
		message = update.EditedMessage
	}

	fmt.Println("here 3", update.PreCheckoutQuery == nil)

	if message != nil {
		if message.SuccessfulPayment != nil && bot.updateHandlers != nil {
			bot.updateHandlers.processSuccessfulPayment(app, &update)
			return
		}

		if message.Chat.Type == "private" {
			state, err := bot.options.stateStorage.UserState(message.Chat.Id)
			if err != nil {
				bot.options.stateStorage.SetUserState(message.Chat.Id, "Welcome")
				bot.updateReplyStateInternalError(update, err)
				return
			}

			if app.MethodByName(state).Kind() == reflect.Invalid {
				bot.updateReplyStateNotExists(update, state)
				return
			}

			apiField := app.Elem().Field(0)
			api.SetRecipientChatId(message.Chat.Id)
			apiField.Set(reflect.ValueOf(api))
			bot.invoke(app, update, state, false)
			return
		}
	}

	fmt.Println("here 4", update.PreCheckoutQuery == nil)

	if bot.updateHandlers == nil {
		return
	}

	fmt.Println("here 5", update.PreCheckoutQuery == nil)

	bot.updateHandlers.processUpdate(app, update)

	return
}

func (bot *Bot) Poll() (err error) {
	gu := bot.api.GetUpdates()

	if !bot.options.getAllUpdates {
		var allowedUpdates = []string{"message", "edited_message"}

		if bot.updateHandlers != nil {
			allowedUpdates = bot.updateHandlers.allowedUpdates()
		}

		gu.SetAllowedUpdates(allowedUpdates)
	}

	go gu.LongPoll()

	for update := range bot.api.Updates() {
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
