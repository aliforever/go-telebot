package app

import (
	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

type App struct {
	go_telegram_bot_api.TelegramBot
	CustomField int
}
