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

func ( actions SayHelloEvents ) OnBotInit( ctx *got.ReactorCtx ) {
    // Things to do when your bot starts for the first time
}

func ( actions SayHelloEvents ) OnText( ctx *got.ReactorCtx ) {
    t := fmt.Sprintf("Hello %s!", msg.Sender.Name)
    ctx.Bot.SendMessage(t, msg.Sender)
}
```

## Conversational Plugin

**Plugin creation**
```go

Colors, _ := got.NewPlugin( got.ConversationalSettings{
    Name: "colors",
    Trigger: "/colors",
    States: ColorsStates,
    StateStartKey: START_LOGIN,
    Events: ColorsEvents{},
    RemindAfter: 2 * 60 * 1000, // reminds every to 2 minutes -> Milliseconds
    ExpireAfter: 10 * 60 * 1000, // 10 minutes -> Milliseconds
    // Storage: YourCustomStorage that follow got.PluginStorage interface
})

```

**Events**
```go
type ColorsEvents struct {}

func ( actions ColorsEvents ) OnBotInit( ctx *ConversationalCtx ) {
}

func ( actions ColorsEvents ) OnSessionStart( ctx *ConversationalCtx ) {
}

func ( actions ColorsEvents ) OnSessionExpired( ctx *ConversationalCtx ) {
}

func ( actions ColorsEvents ) OnSessionRemind( ctx *ConversationalCtx ) {
}

func ( actions ColorsEvents ) OnSessionEnd( ctx *ConversationalCtx ) {
    
}

func ( actions ColorsEvents ) OnAnswer( ctx *ConversationalCtx ) {
}
```

**State struct**
```go
type State struct {
    WaitForAnswer bool
    Finish        bool
    SendQuestion  func( ctx *ConversationalCtx )
    GetNextKey    func( ctx *ConversationalCtx ) (StateKey, bool)
}
```

**State example**

```go

const (
    START_COLORS got.StateKey = iota
    CONFIRM_COLORS
    END_COLORS
)

var ColorsStates got.StatesMap = got.StatesMap{

    START_COLORS: got.State{

        WaitForAnswer: true,

        SendQuestion: func( ctx *got.ConversationalCtx ) {
            ctx.Bot.SendMessage( "What color do you like?", ctx.User )
        },

        GetNextKey: func( ctx *got.ConversationalCtx ) (got.StateKey, bool) {
            return CONFIRM_COLORS, true
        },
    
    },

    CONFIRM_COLORS: got.State{

        WaitForAnswer: true,
        
        SendQuestion: func( ctx *got.ConversationalCtx ) {
            ctx.Bot.SendMessage( "Are you sure? (yes, no)", ctx.User )
        },
    
        GetNextKey: func( ctx *got.ConversationalCtx ) (got.StateKey, bool) {
            
            if answer.Text == "yes" {
                return END_COLORS, true
            }

            if answer.Text == "no" {
                return START_COLORS, true
            }

            ctx.Bot.SendMessage('permitted answers (yes, no)', ctx.User)
            return CONFIRM_COLORS, false
        },

    },

    END_COLORS: got.State{

        Finish: true,
        
        SendQuestion: func( ctx *got.ConversationalCtx ) {
            ctx.Bot.SendMessage( "Thank you! :)", ctx.User )
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
