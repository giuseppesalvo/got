package got

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
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