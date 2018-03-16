package got

import (
	"time"
)

// Types

type ConversationalCtx struct {
	Plugin *ConversationalPlugin
	Bot *Bot
	User *User
	Answer Message
	UserState *UserState
}

type ConversationalEvents interface {
	OnBotInit( ctx *ConversationalCtx )
	OnSessionStart( ctx *ConversationalCtx )
	OnAnswer( ctx *ConversationalCtx )
	OnSessionEnd( ctx *ConversationalCtx )
}

type StateKey int
type StatesMap map[StateKey]State

type ConversationalSettings struct {
	Name          string
	Trigger       string
	States        StatesMap
	StateStartKey StateKey
	Events        ConversationalEvents
	Storage       PluginStorage
}

type ConversationalPlugin struct {
	Name          string
	Trigger       string
	States        StatesMap
	StateStartKey StateKey
	Events        ConversationalEvents
	Storage       PluginStorage
}

type StateCtx struct {
	Bot 	  *Bot
	User 	  *User
	UserState *UserState
	Answer 	  Message
}

type State struct {
	WaitForAnswer bool
	Finish        bool
	SendQuestion  func(ctx *StateCtx)
	GetNextKey    func(ctx *StateCtx) (StateKey, bool)
}

type UserAnswer struct {
	Answer   string
	StateKey StateKey
}

type UserState struct {
	UserId          string
	CurrentStateKey StateKey
	Answers         []UserAnswer
	CreatedAt       time.Time
}

func (state *UserState) getAnswersForStateKey(key StateKey) []UserAnswer {
	answers := []UserAnswer{}

	for _, answ := range state.Answers {
		if answ.StateKey == key {
			answers = append(answers, answ)
		}
	}

	return answers
}

// Functions and Methods

func NewConversationalPlugin(settings ConversationalSettings) *ConversationalPlugin {
	var storage PluginStorage

	if settings.Storage != nil {
		storage = settings.Storage
	} else {
		storage = &MapPluginStorage{
			sessions: make(map[string]*UserState),
		}
	}

	return &ConversationalPlugin{
		Name:          settings.Name,
		Trigger:       settings.Trigger,
		States:        settings.States,
		StateStartKey: settings.StateStartKey,
		Events:        settings.Events,
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

// Run the plugin

func (pl *ConversationalPlugin) run(bot *Bot, msg Message) {

	if !pl.isSessionRunningForUser(msg.Sender) {

		pl.startSession(bot, msg)

	} else {

		// Session is already running

		userState := pl.getUserState(msg.Sender)
		oldState := pl.States[userState.CurrentStateKey]

		if !oldState.Finish {

			pl.goToNextState(bot, msg, userState, oldState)

		} else {

			pl.endSession(bot, msg, userState)

		}
	}
}

// Utils

func (pl *ConversationalPlugin) goToNextState(bot *Bot, msg Message, userState *UserState, state State) {

	ctx := &StateCtx{
		Bot: bot,
		User: msg.Sender,
		UserState: userState,
		Answer: msg,
	}

	newIndex, ok := state.GetNextKey(ctx)

	/*
	 * If the new index is not ok,
	 * means there was an error in the answer, so we will wait for another input
	 * Otherwise, we can add the answer to the userState
	 */

	if ok {

		answer := UserAnswer{
			msg.Text,
			userState.CurrentStateKey,
		}
		userState.Answers = append(userState.Answers, answer)

		pl_ctx := &ConversationalCtx{
			Plugin: pl,
			Bot: bot,
			User: msg.Sender,
			Answer: msg,
			UserState: userState,
		}

		pl.Events.OnAnswer(pl_ctx)

		userState.CurrentStateKey = newIndex
		pl.sendQuestionForUserState(userState, bot, msg)

	}
}

func (pl *ConversationalPlugin) startSession(bot *Bot, msg Message) {
	userState := pl.getUserState(msg.Sender)

	pl_ctx := &ConversationalCtx{
		Plugin: pl,
		Bot: bot,
		User: msg.Sender,
		Answer: msg,
		UserState: userState,
	}

	pl.Events.OnSessionStart(pl_ctx)
	pl.sendQuestionForUserState(userState, bot, msg)
}

func (pl *ConversationalPlugin) endSession(bot *Bot, msg Message, userState *UserState) {

	pl_ctx := &ConversationalCtx{
		Plugin: pl,
		Bot: bot,
		User: msg.Sender,
		Answer: msg,
		UserState: userState,
	}

	pl.Events.OnSessionEnd(pl_ctx)
	pl.Storage.DeleteSessionForUserId(msg.Sender.Id)
}

func (pl *ConversationalPlugin) getUserState(user *User) *UserState {

	userState, ok := pl.Storage.GetSessionFromUserId(user.Id)

	if ok {
		return userState
	} else {
		userState := &UserState{
			user.Id,
			pl.StateStartKey,
			[]UserAnswer{},
			time.Now(),
		}
		pl.Storage.SetSessionForUserId(user.Id, userState)
		return userState
	}
}

func (pl *ConversationalPlugin) isSessionRunningForUser(user *User) bool {
	_, ok := pl.Storage.GetSessionFromUserId(user.Id)
	return ok
}

func (pl *ConversationalPlugin) sendQuestionForUserState(userState *UserState, bot *Bot, currentMsg Message) {

	state := pl.States[userState.CurrentStateKey]

	if state.SendQuestion != nil {

		ctx := &StateCtx{
			Bot: bot,
			User: currentMsg.Sender,
			UserState: userState,
			Answer: currentMsg,
		}

		state.SendQuestion(ctx)
	}

	if !state.WaitForAnswer {
		newMsg := currentMsg
		newMsg.Text = ""
		pl.run(bot, newMsg)
	}
}
