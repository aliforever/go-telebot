package telebot

import (
	"log/slog"
)

type BotOptions struct {
	stateStorage     UserStateStorage
	useStorage       UserStorage
	getAllUpdates    bool
	logger           *slog.Logger
	logRawUpdates    bool
	apiServerURL     *string
	onGetUpdateError func(err error) bool
}

func NewOptions() *BotOptions {
	return &BotOptions{}
}

func (b *BotOptions) SetStateStorage(ss UserStateStorage) *BotOptions {
	b.stateStorage = ss
	return b
}

func (b *BotOptions) SetUserStorage(ss UserStorage) *BotOptions {
	b.useStorage = ss
	return b
}

func (b *BotOptions) GetAllUpdates() *BotOptions {
	b.getAllUpdates = true
	return b
}

func (b *BotOptions) SetLogger(logger *slog.Logger) *BotOptions {
	b.logger = logger
	return b
}

func (b *BotOptions) LogRawUpdates() *BotOptions {
	b.logRawUpdates = true
	return b
}

func (b *BotOptions) SetAPIServerURL(url string) *BotOptions {
	b.apiServerURL = &url
	return b
}

func (b *BotOptions) OnGetUpdateError(f func(err error) bool) *BotOptions {
	b.onGetUpdateError = f
	return b
}
