# Work in progress

# Framework for making bots with go

- Modes 
    - Debug    -> debug your bot from the terminal
    - Telegram -> based on tucnak/telebot.v2
- Plugins
    - Reactor -> Send something after a message
    - Conversational -> It helps you to create conversational commands without loosing the current user state

## Bot

```go

b, err := got.NewBot( got.BotSettings{
    Token: myTelegramToken,
    Mode: got.ModeDebug,
    Plugins: []got.Plugin{ yourplugin },
})

```

## Reactor Plugin

**Plugin creation**
```go

sayhello, _ := got.NewPlugin( got.ConversationalSettings{
    Name: "sayhello",
    Trigger: "/hello", // can be a regex with regexp prefix -> "regexp (hi|hello)"
    Events: SayHelloEvents{},
})
```

**Events**
```go
type SayHelloEvents struct {}

func ( actions SayHelloEvents ) OnBotInit(pl *got.ReactorPlugin, bot *got.Bot) {
    // Things to do when your bot starts for the first time
}

func ( actions SayHelloEvents ) OnText(pl *got.ReactorPlugin, bot *got.Bot, msg got.Message) {
    t := fmt.Sprintf("Hello %s!", msg.Sender.Name)
    bot.SendMessage(t, msg.Sender)
}
```

## Conversational Plugin

**Plugin creation**
```go

Colors, _ := got.NewPlugin( got.ConversationalSettings{
    Name: "black_or_white",
    Trigger: "/black_or_white",
    States: ColorsStates,
    StateStartKey: START_LOGIN,
    Events: ColorsEvents{},
    // Storage: YourCustomStorage that follow got.PluginStorage interface
})

```

**Events**
```go
type ColorsEvents struct {}

func ( actions ColorsEvents ) OnBotInit(pl *got.ConversationalPlugin, bot *got.Bot) {
}

func ( actions ColorsEvents ) OnSessionStart(pl *got.ConversationalPlugin, bot *got.Bot, user *got.User, state *got.UserState) {
}

func ( actions ColorsEvents ) OnSessionEnd(pl *got.ConversationalPlugin, bot *got.Bot, user *got.User, state *got.UserState) {
    
}

func ( actions ColorsEvents ) OnAnswer(pl *got.ConversationalPlugin, bot *got.Bot, user *got.User, answer got.UserAnswer, state *got.UserState) {
}
```

**State struct**
```go
type State struct {
    WaitForAnswer bool
    Finish        bool
    SendQuestion  func(bot *Bot, user *User, state *UserState)
    GetNextKey    func(bot *Bot, user *User, state *UserState, answer Message) (StateKeyType, bool)
}
```

**State example**

```go

const (
    START_COLORS got.StateKeyType = iota
    CONFIRM_COLORS
    END_COLORS
)

var ColorsStates got.StatesMap = got.StatesMap{

    START_COLORS: got.State{

        WaitForAnswer: true,

        SendQuestion: func(bot *got.Bot, user *got.User, state *got.UserState) {
            bot.SendMessage( "What color do you like?", user )
        },

        GetNextKey: func(bot *got.Bot, user *got.User, state *got.UserState, answer got.Message) (got.StateKeyType, bool) {
            return CONFIRM_COLORS, true
        },
    
    },

    CONFIRM_COLORS: got.State{

        WaitForAnswer: true,
        
        SendQuestion: func(bot *got.Bot, user *got.User, state *got.UserState) {
            bot.SendMessage( "Are you sure? (yes, no)", user )
        },
    
        GetNextKey: func(bot *got.Bot, user *got.User, state *got.UserState, answer got.Message) (got.StateKeyType, bool) {
            
            if answer.Text == "yes" {
                return END_COLORS, true
            }

            if answer.Text == "no" {
                return START_COLORS, true
            }

            bot.SendMessage('permitted answers (yes, no)')
            return CONFIRM_COLORS, false
        },

    },

    END_COLORS: got.State{

        Finish: true,
        
        SendQuestion: func(bot *got.Bot, user *got.User, state *got.UserState) {
            bot.SendMessage( "Thank you! :)", user )
        },
    
    },
}

```

**Storage Interface**

```go

type PluginStorage interface {
    GetSessionFromUserId(id string) (*UserState, bool)
    SetSessionForUserId(id string, state *UserState)
    DeleteSessionForUserId(id string)
}

```

**In memory storage example**

```go

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

```

**Keyboard Markup**

```go
markup := &got.ReplyMarkup{
    ReplyKeyboard: [][]got.ReplyButton{
        []got.ReplyButton{
            got.ReplyButton{Text: "yes"},
            got.ReplyButton{Text: "no"},
        },
    },
})

bot.SendMessage("Are you sure", sender, markup)
```
