package main

import (
	"fmt"
	"os"
	"l4tt/gobrute/modules/scraper"
	"l4tt/gobrute/modules/error"
	"l4tt/gobrute/modules/logger"
	"l4tt/gobrute/modules/check_balance"
)


func main() {
	args := os.Args
	args_list := []string{"--help", "--target", "--passwordlist"}
	target_list := []string{"checkbalance", "scrape"}


	if len(args) < 2 {
		fmt.Println("Usage: gobrute --help")
		return
	}

	if args[1] == args_list[0] {
		fmt.Println(`
		╔══════════════════════════════════════════════╗
		║ Usage: gobrute 1.0.0                         ║
		║                                              ║
		║ Options:                                     ║
		║   --target [checkbalance, scrape]            ║
		║   --passwordlist [dictionary]                ║
		╚══════════════════════════════════════════════╝
		`)
	}
	if args[1] == args_list[1] {
		if args[2] == target_list[0] {
			logger.Log("Target: " + args[2], false)
			check_balance.CheckBalance()
		} else if args[2] == target_list[1] {
			if len(args) > 3 {
				scraper.Scrape(args[2], args[4])
			} else {
				error.ErrorHandler("Password list not specified")
			}
		} else {
			error.ErrorHandler("Invalid target specified")
		}
	}
}


