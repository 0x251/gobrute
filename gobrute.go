package main

import (
	"flag"
	"fmt"
	"l4tt/gobrute/modules/check_balance"
	"l4tt/gobrute/modules/error"
	"l4tt/gobrute/modules/generator"
	"l4tt/gobrute/modules/logger"
	"l4tt/gobrute/modules/scraper"
)

func main() {
	target := flag.String("target", "", "Specify the target [checkbalance, scrape, localcheck, generate]")
	passwordList := flag.String("passwordlist", "", "Specify the password list file")
	threads := flag.Int("threads", 10, "Number of concurrent threads (0-150)")
	addressList := flag.String("addresslist", "", "Specify the address list file")
	help := flag.Bool("help", false, "Display help")
	generate := flag.Bool("generate", false, "Generate a password list")
	random := flag.Bool("random", false, "Generate a random password list")
	multiWord := flag.Bool("multi-word", false, "Generate a multi-word password list")
	count := flag.Int("count", 100, "Number of passwords to generate")
	filename := flag.String("filename", "passwords.txt", "Filename to save the generated passwords")

	flag.Parse()

	if *random || *multiWord {
		*generate = true
	}

	if *help || *target == "" {
		fmt.Println(`
		╔════════════════════════════════════════════════╗
		║ Usage: gobrute 1.0.1                            ║
		║                                                ║
		║ Options:                                       ║
		║   --target [checkbalance, scrape, localcheck, generate]  ║
		║   --passwordlist [dictionary]                  ║
		║   --threads [number]                           ║
		║   --addresslist [address list file]           ║
		║   --generate                                   ║
		║   --random                                     ║
		║   --multi-word                                 ║
		║   --count [number]                             ║
		║   --filename [filename]                        ║
		╚════════════════════════════════════════════════╝
		`)
		return
	}

	switch *target {
	case "checkbalance":
		if *threads <= 0 {
			*threads = 10
		}

		if *threads > 150 {
			logger.Log("Too many concurrent threads, please use a lower number, max 150", true)
			return
		}

		check_balance.CheckBalance(*threads)
	case "scrape":
		if *passwordList == "" {
			error.ErrorHandler("Password list not specified, please use --passwordlist [dictionary].txt")
			return
		}
		scraper.Scrape(*target, *passwordList)
	case "localcheck":
		if *addressList == "" {
			error.ErrorHandler("Address list not specified, please use --addresslist [address list file].txt")
			return
		}
		scraper.LocalCheck(*addressList)
	case "generate":
		if !*generate {
			error.ErrorHandler("Generate flag not set. Use --random or --multi-word to generate passwords.")
			return
		}

		if *random {
			generator.GeneratePasswords(true, false, *count, *filename)
		} else if *multiWord {
			generator.GeneratePasswords(false, true, *count, *filename)
		} else {
			error.ErrorHandler("Specify generation type: --random or --multi-word")
		}
	default:
		error.ErrorHandler("Invalid target specified, please use --target [checkbalance, scrape, localcheck, generate]")
	}
}
