package got

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"time"
	"strings"
)

type BotMode int

const (
	ModeDebug BotMode = iota
	ModeTelegram
)

type User struct {
	Id       string
	Name     string
	Telegram *tb.User
}

type Chat struct {
	Id       string
	Telegram *tb.Chat
}

type Message struct {
	Id       string
	Text     string
	Sender   *User
	Chat     *Chat
	Date     time.Time
	Telegram *tb.Message
}

type BotSettings struct {
	Token   string
	Mode    BotMode
	Plugins []Plugin
}

type Bot struct {
	Settings BotSettings
	Telegram *tb.Bot
}

func NewBot(settings BotSettings) (*Bot, error) {
	return &Bot{Settings: settings}, nil
}

// Bot Methods

func normalizeTelegramReplyMarkup(markup *ReplyMarkup) *tb.ReplyMarkup {

	tb_markup := &tb.ReplyMarkup{}

	if markup.ReplyKeyboard != nil {

		keyboard := [][]tb.ReplyButton{}

		for _, row := range markup.ReplyKeyboard {
			tb_row := []tb.ReplyButton{}
			for _, col := range row {
				tb_row = append(tb_row, tb.ReplyButton{
					Text: col.Text,
				})
			}
			keyboard = append(keyboard, tb_row)
		}

		tb_markup.ReplyKeyboard = keyboard
	}

	if markup.ReplyKeyboardRemove {
		tb_markup.ReplyKeyboardRemove = markup.ReplyKeyboardRemove
	}

	fmt.Printf("\n%+v\n%+v\n", markup, tb_markup)

	return tb_markup
}

func normalizeDebugReplyMarkup(markup *ReplyMarkup) string {

	if markup.ReplyKeyboardRemove {
		return ""
	}

	tb_markup := ""

	if markup.ReplyKeyboard != nil {

		keyboard := []string{}

		for _, row := range markup.ReplyKeyboard {
			tb_row := []string{}
			for _, col := range row {
				tb_row = append(tb_row, col.Text)
			}
			if len(tb_row) > 0 {
				keyboard = append( keyboard, "[" + strings.Join(tb_row, ", ") + "]" )
			}
		}

		tb_markup = strings.Join(keyboard, "\n")
	}

	return tb_markup
}

func normalizeTelegramOptions(options interface{}) interface{} {
	switch normalized := options.(type) {
	case *ReplyMarkup:
		return normalizeTelegramReplyMarkup(normalized)
	default:
		panic("Invalid options")
	}
}

func normalizeDebugOptions(options interface{}) interface{} {
	switch normalized := options.(type) {
	case *ReplyMarkup:
		return normalizeDebugReplyMarkup(normalized)
	default:
		panic("Invalid options")
	}
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

// Utils

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
