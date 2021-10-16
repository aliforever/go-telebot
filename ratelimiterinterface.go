package telebot

import go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"

type RateLimiter interface {
	ShouldLimitUpdate(*go_telegram_bot_api.Update) bool
}
