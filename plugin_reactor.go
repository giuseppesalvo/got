package got

// Types

type ReactorEvents interface {
	OnBotInit(pl *ReactorPlugin, bot *Bot)
	OnText(pl *ReactorPlugin, bot *Bot, msg Message)
}

type ReactorSettings struct {
	Name    string
	Trigger string
	Events  ReactorEvents
}

type ReactorPlugin struct {
	Name    string
	Trigger string
	Events  ReactorEvents
}

// Functions and Methods

func NewReactorPlugin(settings ReactorSettings) *ReactorPlugin {
	return &ReactorPlugin{
		Name:    settings.Name,
		Trigger: settings.Trigger,
		Events:  settings.Events,
	}
}

// Methods that implements the Plugin interface

func (pl *ReactorPlugin) onInit(bot *Bot) {
	pl.Events.OnBotInit(pl, bot)
}

func (pl *ReactorPlugin) onText(bot *Bot, msg Message) {

	triggered := checkTriggerInStr(pl.Trigger, msg.Text)

	if triggered {
		pl.run(bot, msg)
	}
}

// Run the plugin

func (pl *ReactorPlugin) run(bot *Bot, msg Message) {
	pl.Events.OnText(pl, bot, msg)
}
