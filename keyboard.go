package got

type ReplyMarkup struct {
	ReplyKeyboardRemove bool
	ReplyKeyboard       [][]ReplyButton
}

type ReplyButton struct {
	Text string
}
