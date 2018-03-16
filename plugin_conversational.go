package got

import (
	"time"
)

// Types

type ConversationalEvents interface {
	OnBotInit(pl *ConversationalPlugin, bot *Bot)
	OnSessionStart(pl *ConversationalPlugin, bot *Bot, user *User, state *UserState)
	OnAnswer(pl *ConversationalPlugin, bot *Bot, user *User, answer UserAnswer, state *UserState)
	OnSessionEnd(pl *ConversationalPlugin, bot *Bot, user *User, state *UserState)
}

type StateKeyType int
type StatesMap map[StateKeyType]State

type ConversationalSettings struct {
	Name          string
	Trigger       string
	States        StatesMap
	StateStartKey StateKeyType
	Events        ConversationalEvents
	Storage       PluginStorage
}

type ConversationalPlugin struct {
	Name          string
	Trigger       string
	States        StatesMap
	StateStartKey StateKeyType
	Events        ConversationalEvents
	Storage       PluginStorage
}

type State struct {
	WaitForAnswer bool
	Finish        bool
	SendQuestion  func(bot *Bot, user *User, state *UserState)
	GetNextKey    func(bot *Bot, user *User, state *UserState, answer Message) (StateKeyType, bool)
}

type UserAnswer struct {
	Answer   string
	StateKey StateKeyType
}

type UserState struct {
	UserId          string
	CurrentStateKey StateKeyType
	Answers         []UserAnswer
	CreatedAt       time.Time
}

func (state *UserState) getAnswersForStateKey(key StateKeyType) []UserAnswer {
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
	pl.Events.OnBotInit(pl, bot)
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

	newIndex, ok := state.GetNextKey(bot, msg.Sender, userState, msg)

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

		pl.Events.OnAnswer(pl, bot, msg.Sender, answer, userState)

		userState.CurrentStateKey = newIndex
		pl.sendQuestionForUserState(userState, bot, msg)

	}
}

func (pl *ConversationalPlugin) startSession(bot *Bot, msg Message) {
	userState := pl.getUserState(msg.Sender)
	pl.Events.OnSessionStart(pl, bot, msg.Sender, userState)
	pl.sendQuestionForUserState(userState, bot, msg)
}

func (pl *ConversationalPlugin) endSession(bot *Bot, msg Message, userState *UserState) {
	pl.Events.OnSessionEnd(pl, bot, msg.Sender, userState)
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
		state.SendQuestion(bot, currentMsg.Sender, userState)
	}

	if !state.WaitForAnswer {
		newMsg := currentMsg
		newMsg.Text = ""
		pl.run(bot, newMsg)
	}
}
