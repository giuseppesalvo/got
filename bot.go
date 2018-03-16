package got

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"time"
)

type BotMode int

const (
	ModeDebug BotMode = iota
	ModeTelegram
)

type BotSettings struct {
	Token   string
	Mode    BotMode
	Plugins []Plugin
}

type Bot struct {
	Settings BotSettings
	Telegram *tb.Bot
}

/**
 * Bot Methods
 * 
 */

func NewBot(settings BotSettings) (*Bot, error) {
	return &Bot{Settings: settings}, nil
}

func (bot *Bot) SendMessage(msg string, sender *User, options ...interface{}) {
	if bot.Settings.Mode == ModeDebug {
		fmt.Println("🤖: " + msg)

		if len(options) > 0 {
			markup := normalizeDebugOptions(options[0]).(string)
			if len(markup) > 0 {
				fmt.Println(markup)
			}
		}

	} else if bot.Settings.Mode == ModeTelegram {
		if len(options) > 0 {
			bot.Telegram.Send(sender.Telegram, msg, normalizeTelegramOptions(options[0]))
		} else {
			bot.Telegram.Send(sender.Telegram, msg)
		}
	}
}

// Start functions

func (bot *Bot) Start() {

	if bot.Settings.Mode == ModeDebug {

		bot.startDebug()

	} else if bot.Settings.Mode == ModeTelegram {

		bot.startTelegram()

	}
}

// Debug in command line

func (bot *Bot) startDebug() {

	bot.pluginsOnInit()

	for {
		fmt.Print("→ ")
		res := readConsoleLine()

		if res == "exit()" {
			break
		}

		user_id := "terminal_user"

		if regMatchStr("[A-Z]{2,2}:\\s.+", res) {
			user_id = res[0:2]
			res = res[4:]
		}

		user := &User{
			Id:   user_id,
			Name: user_id,
		}

		chat := &Chat{
			Id: "terminal_chat",
		}

		msg := Message{
			Id:     "0",
			Text:   res,
			Sender: user,
			Chat:   chat,
			Date:   time.Now(),
		}

		bot.pluginsOnText(msg)
	}
}

// Telegram

func (bot *Bot) startTelegram() {

	b, err := tb.NewBot(tb.Settings{
		Token:  bot.Settings.Token,
		Poller: &tb.LongPoller{Timeout: 5 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	bot.pluginsOnInit()

	bot.Telegram = b

	b.Handle(tb.OnText, func(m *tb.Message) {

		user := &User{
			Id:       strconv.Itoa(int(m.Sender.ID)),
			Name:     m.Sender.FirstName,
			Telegram: m.Sender,
		}

		chat := &Chat{
			Id:       strconv.Itoa(int(m.Chat.ID)),
			Telegram: m.Chat,
		}

		msg := Message{
			Id:       strconv.Itoa(int(m.ID)),
			Text:     m.Text,
			Sender:   user,
			Chat:     chat,
			Date:     m.Time(),
			Telegram: m,
		}

		bot.pluginsOnText(msg)
	})

	b.Start()
}

/**
 * Utils
 * 
 */

func (bot *Bot) arePluginsOk() bool {
	return bot.Settings.Plugins != nil && len(bot.Settings.Plugins) > 0
}

func (bot *Bot) pluginsOnInit() {
	if bot.arePluginsOk() {
		for _, plugin := range bot.Settings.Plugins {
			plugin.(Plugin).onInit(bot)
		}
	}
}

func (bot *Bot) pluginsOnText(msg Message) {
	if bot.arePluginsOk() {
		for _, plugin := range bot.Settings.Plugins {
			plugin.(Plugin).onText(bot, msg)
		}
	}
}

func normalizeTelegramOptions(options interface{}) interface{} {
	switch typed := options.(type) {
	case *ReplyMarkup:
		return typed.normalizeForTelegram()
	default:
		panic("Invalid options")
	}
}

func normalizeDebugOptions(options interface{}) interface{} {
	switch typed := options.(type) {
	case *ReplyMarkup:
		return typed.normalizeForDebug()
	default:
		panic("Invalid options")
	}
}