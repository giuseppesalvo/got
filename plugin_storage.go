package got

/**
 * Customizable PluginStorage for retrieving user sessions in conversational plugins
 */

type PluginStorage interface {
	GetSessionByUserId(id string) (*Session, bool)
	SetSessionForUserId(id string, state *Session)
	DeleteSessionForUserId(id string)
}

/**
 * In memory storage with map[string]*Session
 * Default bot storage
 */

type MapPluginStorage struct {
	sessions map[string]*Session
}

func (storage *MapPluginStorage) GetSessionByUserId(id string) (*Session, bool) {
	state, ok := storage.sessions[id]
	return state, ok
}

func (storage *MapPluginStorage) SetSessionForUserId(id string, state *Session) {
	storage.sessions[id] = state
}

func (storage *MapPluginStorage) DeleteSessionForUserId(id string) {
	delete(storage.sessions, id)
}
