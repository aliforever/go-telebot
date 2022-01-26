package app

import "github.com/aliforever/go-telegram-bot-api"

type App struct {
	tgbotapi.TelegramBot
	CustomField int
}
