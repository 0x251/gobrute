package logger

import (
	"fmt"
	"time"
)

func Log(message string, warning bool) {
	if warning {
		fmt.Printf("\033[33m[%s] GoBrute Warning: %s\033[0m\n", time.Now().Format("2006-01-02 15:04:05"), message)
	} else {
		fmt.Printf("\033[37m[%s] GoBrute: %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	}
}

func Status(message string) {
	fmt.Printf("\r\033[K%s", message)
	fmt.Print("\r")
}
