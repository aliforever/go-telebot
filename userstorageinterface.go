package telebot

import (
	"github.com/aliforever/go-telegram-bot-api/structs"
)

type UserStorage interface {
	Store(user *structs.User) error
}
