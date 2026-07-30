package types

import (
	"fmt"
)

type COLOR string

const (
	Reset  COLOR = "\033[0m"
	Red    COLOR = "\033[31m"
	Green  COLOR = "\033[32m"
	Yellow COLOR = "\033[33m"
	Blue   COLOR = "\033[34m"
)

func TermSetColor(c COLOR) {
	fmt.Printf("%s", c)
}

func TermResetColor() {
	fmt.Printf("%s", Reset)
}
