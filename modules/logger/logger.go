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

func UpdateStatus(passwords int, left int) {
	fmt.Printf("\033[34m[ Passwords: %d ] %% [ Left: %d ]\033[0m\033[K", passwords, left)
}