package got

type PluginType int

type Plugin interface {
	onText(bot *Bot, msg Message)
	onInit(bot *Bot)
}

func NewPlugin(settings interface{}) (Plugin, error) {
	switch setts := settings.(type) {
	case ConversationalSettings:
		return NewConversationalPlugin(setts), nil
	case ReactorSettings:
		return NewReactorPlugin(setts), nil
	default:
		panic("invalid settings")
	}
}
