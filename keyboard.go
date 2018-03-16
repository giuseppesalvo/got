package got

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

type ReplyMarkup struct {
	ReplyKeyboardRemove bool
	ReplyKeyboard       [][]ReplyButton
}

type ReplyButton struct {
	Text string
}

func (markup *ReplyMarkup) normalizeForTelegram() *tb.ReplyMarkup {

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

	return tb_markup
}

func (markup *ReplyMarkup) normalizeForDebug() string {

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
				keyboard = append(keyboard, "["+strings.Join(tb_row, ", ")+"]")
			}
		}

		tb_markup = strings.Join(keyboard, "\n")
	}

	return tb_markup
}
