package got

// Types

type ReactorCtx struct {
	Plugin  *ReactorPlugin
	Bot     *Bot
	Message Message
	User    *User
}

type ReactorEvents interface {
	OnBotInit( ctx *ReactorCtx )
	OnText( ctx *ReactorCtx )
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
	ctx := &ReactorCtx{
		Plugin: pl,
		Bot: bot,
	}

	pl.Events.OnBotInit(ctx)
}

func (pl *ReactorPlugin) onText(bot *Bot, msg Message) {

	triggered := checkTriggerInStr(pl.Trigger, msg.Text)

	if triggered {
		pl.run(bot, msg)
	}
}

// Run the plugin

func (pl *ReactorPlugin) run(bot *Bot, msg Message) {

	ctx := &ReactorCtx{
		Plugin: pl,
		Bot: bot,
		Message: msg,
		User: msg.Sender,
	}

	pl.Events.OnText(ctx)
}
