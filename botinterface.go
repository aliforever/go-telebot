package telebot

type UserStateStorage interface {
	UserState(userId int64) (string, error)
	SetUserState(userId int64, state string) error
}
