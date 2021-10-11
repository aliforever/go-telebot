package telebot

import "sync"

type stateStorage struct {
	m            sync.Mutex
	data         map[int64]string
	defaultState string
}

func newStateStorage() *stateStorage {
	return &stateStorage{
		data:         map[int64]string{},
		defaultState: "Welcome",
	}
}

func (ss *stateStorage) setDefaultState(defaultState string) *stateStorage {
	ss.defaultState = defaultState
	return ss
}

func (ss *stateStorage) SetUserState(userId int64, state string) error {
	ss.m.Lock()
	defer ss.m.Unlock()
	ss.data[userId] = state
	return nil
}

func (ss *stateStorage) UserState(userId int64) (state string, err error) {
	var ok bool
	if state, ok = ss.data[userId]; ok {
		return
	}
	state = ss.defaultState
	return
}
