package bot

import "fmt"

type BotError struct {
	description string
}

func (err *BotError) Error() string {
	return fmt.Sprintf(err.description)
}
