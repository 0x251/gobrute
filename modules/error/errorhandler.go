package error

import (
	"fmt"
	"os"
	"time"
)

func ErrorHandler(errorMsg string) {
	fmt.Printf("[%s] GoBrute Error: %s\n", time.Now().Format("2006-01-02 15:04:05"), errorMsg)
	os.Exit(1)
}
