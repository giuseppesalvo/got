package got

import (
	"time"
	"github.com/giuseppesalvo/tm"
)

// Types

type ConversationalCtx struct {
	Plugin 	*ConversationalPlugin
	Bot    	*Bot
	User   	*User
	Answer  Message
	Session *Session
}

type ConversationalEvents interface {
	OnAnswer( ctx *ConversationalCtx )
	OnBotInit( ctx *ConversationalCtx )
	OnSessionEnd( ctx *ConversationalCtx )
	OnSessionStart( ctx *ConversationalCtx )
	OnSessionRemind( ctx *ConversationalCtx )
	OnSessionExpired( ctx *ConversationalCtx )
}

type ConversationalSettings struct {
	Name           string
	Trigger        string
	States         States
	Events         ConversationalEvents
	Storage        PluginStorage
	RemindEvery    time.Duration
	ExpireAfter    time.Duration
}

type ConversationalPlugin struct {
	Name           string
	Trigger        string
	States         States
	Events         ConversationalEvents
	Storage        PluginStorage
	RemindEvery    time.Duration
	ExpireAfter    time.Duration
}

func NewConversationalPlugin(settings ConversationalSettings) *ConversationalPlugin {
	var storage PluginStorage

	if settings.Storage != nil {
		storage = settings.Storage
	} else {
		storage = &MapPluginStorage{
			sessions: make(map[string]*Session),
		}
	}

	return &ConversationalPlugin{
		Name:          settings.Name,
		Trigger:       settings.Trigger,
		States:        settings.States,
		Events:        settings.Events,
		RemindEvery:   settings.RemindEvery,
		ExpireAfter:   settings.ExpireAfter,
		Storage:       storage,
	}
}

// Methods that implements the Plugin interface

func (pl *ConversationalPlugin) onInit(bot *Bot) {
	ctx := &ConversationalCtx{
		Plugin: pl,
		Bot: bot,
	}
	pl.Events.OnBotInit(ctx)
}

func (pl *ConversationalPlugin) onText(bot *Bot, msg Message) {

	triggered := checkTriggerInStr(pl.Trigger, msg.Text)

	if triggered || pl.isSessionRunningForUser(msg.Sender) {
		pl.run(bot, msg)
	}
}

/**
 * Run the plugin
 * This function is triggered from onText only if the message contains the plugin trigger,
 * or the user session in currently running
 */

func (pl *ConversationalPlugin) run(bot *Bot, msg Message) {

	session := pl.getSession(msg.Sender)
	pl.setRemindIntervalToSession(session, bot)
	pl.setExpireTimeoutToSession(session, bot)

	if !pl.isSessionRunningForUser(msg.Sender) {
		pl.startSession(bot, msg)
	} else {
		pl.sendQuestionForSession(session, bot, msg)
	}
}

/**
 * Timeout methods
 *
 */

func ( pl *ConversationalPlugin ) clearExpireTimeoutToSession( session *Session, bot *Bot ) {
	tm.ClearTimeout(session.ExpireTimeout)
}

func ( pl *ConversationalPlugin ) setExpireTimeoutToSession( session *Session, bot *Bot ) {
	if pl.ExpireAfter > 0 {

		pl.clearExpireTimeoutToSession(session, bot)

		session.ExpireTimeout = tm.SetTimeout(func () {

			pl.clearRemindIntervalToSession(session, bot)

			pl.Storage.DeleteSessionForUserId(session.UserId)

			pl_ctx := &ConversationalCtx{
				Plugin: pl,
				Bot: bot,
				Session: session,
			}

			pl.Events.OnSessionExpired(pl_ctx)

		}, pl.ExpireAfter)
	}
} 

/**
 * Interval methods
 *
 */

func ( pl *ConversationalPlugin ) clearRemindIntervalToSession( session *Session, bot *Bot ) {
	tm.ClearInterval(session.RemindInterval)
}

func ( pl *ConversationalPlugin ) setRemindIntervalToSession( session *Session, bot *Bot ) {
	if pl.RemindEvery > 0 {

		pl.clearRemindIntervalToSession(session, bot)

		session.RemindInterval = tm.SetInterval(func () {

			pl_ctx := &ConversationalCtx{
				Plugin: pl,
				Bot: bot,
				Session: session,
			}

			pl.Events.OnSessionRemind(pl_ctx)

		}, pl.RemindEvery)
	}
}

/**
 * Session methods
 *
 */

func ( pl *ConversationalPlugin ) RepeatSessionFromCtx( ctx *ConversationalCtx ) {
	session := pl.getSession(ctx.User)
	pl.sendQuestionForSession(session, ctx.Bot, ctx.Answer)
}

func (pl *ConversationalPlugin) startSession(bot *Bot, msg Message) {
	session := pl.getSession(msg.Sender)
	session.Running = true

	pl_ctx := &ConversationalCtx{
		Plugin: pl,
		Bot: bot,
		User: msg.Sender,
		Answer: msg,
		Session: session,
	}

	pl.Events.OnSessionStart(pl_ctx)
	pl.sendQuestionForSession(session, bot, msg)
}

func (pl *ConversationalPlugin) endSession(bot *Bot, msg Message, session *Session) {

	pl.clearRemindIntervalToSession(session, bot)
	pl.clearExpireTimeoutToSession(session, bot)

	pl_ctx := &ConversationalCtx{
		Plugin: pl,
		Bot: bot,
		User: msg.Sender,
		Answer: msg,
		Session: session,
	}

	pl.Events.OnSessionEnd(pl_ctx)
	pl.Storage.DeleteSessionForUserId(msg.Sender.Id)
}

func (pl *ConversationalPlugin) getSession(user *User) *Session {

	session, ok := pl.Storage.GetSessionByUserId(user.Id)

	if ok {
		return session
	} else {
		session := &Session{
			UserId: user.Id,
			StateIndex: 0,
			Answers: []UserAnswer{},
			CreatedAt: time.Now(),
			Cronology: []Message{},
		}
		pl.Storage.SetSessionForUserId(user.Id, session)
		return session
	}
}

func (pl *ConversationalPlugin) isSessionRunningForUser(user *User) bool {
	session, ok := pl.Storage.GetSessionByUserId(user.Id)
	if ok && session.Running {
		return true
	}
	return false
}

func (pl *ConversationalPlugin) sendQuestionForSession(session *Session, bot *Bot, currentMsg Message) {

	state := pl.States[session.StateIndex]

	ctx := &ConversationalCtx{
		Plugin: pl,
		Bot: bot,
		User: currentMsg.Sender,
		Session: session,
		Answer: currentMsg,
	}

	state(ctx)
}
