package scraper

import (
	"encoding/json"
	"fmt"
	"l4tt/gobrute/modules/filemanager"
	"l4tt/gobrute/modules/logger"
	"time"
)

func LocalCheck(address_list string) {
	startTime := time.Now()

	data, err := filemanager.ReadFile(address_list)
	if err != nil {
		logger.Log("Error reading address list: "+err.Error(), true)
		return
	}

	privateKeys, err := filemanager.ReadJSON("results/private_keys.json")
	if err != nil {
		logger.Log("Error reading private keys: "+err.Error(), true)
		return
	}

	addressMap := make(map[string]struct{})
	for _, addr := range data {
		addressMap[addr] = struct{}{}
	}

	if len(privateKeys) == 0 {
		logger.Log("private_keys.json is empty", true)
		return
	}

	keysData, ok := privateKeys[0].(map[string]interface{})
	if !ok {
		logger.Log("Invalid format for private_keys.json", true)
		return
	}

	keys, ok := keysData["keys"].(map[string]interface{})
	if !ok {
		logger.Log("Invalid 'keys' format in private_keys.json", true)
		return
	}

	totalToCheck := len(keys)
	var processed int
	var localCheckList []map[string]string

	for _, key := range keys {
		keyInfo, ok := key.(map[string]interface{})
		if !ok {
			continue
		}

		address, aOk := keyInfo["address"].(string)
		privateKey, pOk := keyInfo["private_key"].(string)
		publicKey, pubOk := keyInfo["public_key"].(string)

		if aOk && pOk && pubOk {
			if _, exists := addressMap[address]; exists {
				localCheckList = append(localCheckList, map[string]string{
					"address":     address,
					"private_key": privateKey,
					"public_key":  publicKey,
				})
			}
		}
		processed++

		remaining := totalToCheck - processed
		elapsed := time.Since(startTime)
		var ttc time.Duration
		if processed > 0 {
			avgTime := elapsed / time.Duration(processed)
			ttc = avgTime * time.Duration(remaining)
		} else {
			ttc = 0
		}

		logger.Status(fmt.Sprintf(
			"Remaining: [\033[34m%d\033[0m], Time to completion: [\033[34m%s\033[0m] Addresses: [\033[34m%d\033[0m]",
			remaining,
			ttc.Round(time.Second),
			len(localCheckList),
		))
	}

	jsonData, err := json.MarshalIndent(localCheckList, "", "  ")
	if err != nil {
		logger.Log("Error marshaling local check data: "+err.Error(), true)
		return
	}

	err = filemanager.WriteFile("results/localcheck.json", []string{string(jsonData)})
	if err != nil {
		logger.Log("Error writing to localcheck.json: "+err.Error(), true)
		return
	}

	elapsed := time.Since(startTime)
	remaining := totalToCheck - processed
	if remaining < 0 {
		remaining = 0
	}
	logger.Log(fmt.Sprintf(
		"Local check completed in [\033[34m%s\033[0m], saved to localcheck.json total [\033[34m%d\033[0m] addresses. Total to check: [\033[34m%d\033[0m], Remaining: [\033[34m%d\033[0m]",
		elapsed.Round(time.Second),
		len(localCheckList),
		totalToCheck,
		remaining,
	), false)
}
