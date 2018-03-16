package bot

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func readConsoleLine() string {
	reader := bufio.NewReader(os.Stdin)
	bytes, _, _ := reader.ReadLine()
	return string(bytes)
}

func regMatchStr(reg string, str string) bool {
	matched, err := regexp.MatchString(reg, str)
	if err != nil {
		return false
	} else {
		return matched
	}
}

func checkTriggerInStr(trigger string, msg string) bool {
	if strings.HasPrefix(trigger, "regexp ") {
		return regMatchStr(trigger[7:], msg)
	} else {
		return msg == trigger
	}
}
