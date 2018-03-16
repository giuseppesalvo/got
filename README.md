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

func ( actions SayHelloEvents ) OnBotInit(pl *got.ReactorPlugin, b *got.Bot) {
    // Things to do when your bot starts for the first time
}

func ( actions SayHelloEvents ) OnText(pl *got.ReactorPlugin, b *got.Bot, msg got.Message) {
    t := fmt.Sprintf("Hello %s!", msg.Sender.Name)
    b.SendMessage(t, msg.Sender)
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

**States**
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
            got.SendMessage( "What color do you like?", user )
        },

        GetNextKey: func(bot *got.Bot, user *got.User, state *got.UserState, answer got.Message) (got.StateKeyType, bool) {
            return CONFIRM_COLORS, true
        },
    
    },

    CONFIRM_COLORS: got.State{

        WaitForAnswer: true,
        
        SendQuestion: func(bot *got.Bot, user *got.User, state *got.UserState) {
            got.SendMessage( "Are you sure? ( yes, no )", user )
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
            got.SendMessage( "Thank you! :)", user )
        },
    
    },
}

```
