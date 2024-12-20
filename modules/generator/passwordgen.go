package generator

import (
	"fmt"
	"l4tt/gobrute/modules/filemanager"
	"l4tt/gobrute/modules/logger"
	"math/rand"
	"time"
)

func GeneratePasswords(random bool, multi_word bool, count int, filename string) {
	words, err := filemanager.ReadFile("results/bip39_words.txt")
	filemanager.DeleteFile(filename)
	if err != nil {
		logger.Log("Error reading words file: "+err.Error(), true)
		return
	}

	var passwords []string

	rand.Seed(time.Now().UnixNano())

	currentIndex := 0
	startTime := time.Now()

	for i := 0; i < count; i++ {
		var password string
		if multi_word {
			if len(words) < 2 {
				logger.Log("Not enough words to generate a multi-word password.", true)
				return
			}
			word1 := words[rand.Intn(len(words))]
			word2 := words[rand.Intn(len(words))]
			password = word1 + " " + word2
		} else {
			if len(words) < 1 {
				logger.Log("No words available to generate a password.", true)
				return
			}
			if random {
				password = words[rand.Intn(len(words))]
			} else {
				password = words[currentIndex%len(words)]
				currentIndex++
			}
		}
		passwords = append(passwords, password)

		generated := i + 1
		remaining := count - generated
		elapsed := time.Since(startTime)
		var eta time.Duration
		if generated > 0 {
			avgTimePer := elapsed / time.Duration(generated)
			eta = avgTimePer * time.Duration(remaining)
		} else {
			eta = 0
		}

		logger.Status(fmt.Sprintf("Generated: %d, Remaining: %d, ETA: %s", generated, remaining, eta.Round(time.Second)))
	}

	err = filemanager.WritePassword(filename, passwords)
	if err != nil {
		logger.Log("Error writing passwords to file: "+err.Error(), true)
		return
	}

	logger.Log(fmt.Sprintf("Successfully generated %d passwords and saved to %s", count, filename), false)
}
