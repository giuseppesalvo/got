package got

/**
 * Customizable PluginStorage for retrieving user sessions in conversational plugins
 */

type PluginStorage interface {
	GetSessionFromUserId(id string) (*UserState, bool)
	SetSessionForUserId(id string, state *UserState)
	DeleteSessionForUserId(id string)
}

/**
 * In memory storage with map[string]*UserState
 * Default bot storage
 */

type MapPluginStorage struct {
	sessions map[string]*UserState
}

func (storage *MapPluginStorage) GetSessionFromUserId(id string) (*UserState, bool) {
	state, ok := storage.sessions[id]
	return state, ok
}

func (storage *MapPluginStorage) SetSessionForUserId(id string, state *UserState) {
	storage.sessions[id] = state
}

func (storage *MapPluginStorage) DeleteSessionForUserId(id string) {
	delete(storage.sessions, id)
}
