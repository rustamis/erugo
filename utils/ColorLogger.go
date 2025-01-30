package utils

import (
	"fmt"
	"log"
)

//Provide a log function that allows for colorized output

const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorReset  = "\033[0m"
)

func Log(message string, color string, args ...any) {
	log.Printf("%s%s%s", color, fmt.Sprintf(message, args...), ColorReset)
}
