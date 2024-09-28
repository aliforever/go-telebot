package telebot

type BotOptions struct {
	stateStorage  UserStateStorage
	useStorage    UserStorage
	getAllUpdates bool
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
