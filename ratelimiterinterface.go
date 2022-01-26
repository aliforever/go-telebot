package telebot

import "github.com/aliforever/go-telegram-bot-api"

type RateLimiter interface {
	ShouldLimitUpdate(*tgbotapi.Update) bool
}
