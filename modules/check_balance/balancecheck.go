package check_balance

import (
	"encoding/json"
	"fmt"
	"l4tt/gobrute/modules/filemanager"
	"l4tt/gobrute/modules/logger"
	"net/http"
	"sync"
	"time"
)

var mu sync.Mutex

func CheckBalance() {
	jsonData, err := filemanager.ReadJSON("private_keys.json")
	if err != nil {
		logger.Log("Failed to read Addresses, please make sure you have scraped and a file called private_keys.json exists", true)
		return
	}

	totalAddresses := 0
	addressWithTxs := 0
	var wg sync.WaitGroup

	for _, item := range jsonData {
		if data, ok := item.(map[string]interface{}); ok {
			totalAddresses = int(data["count"].(float64))
			if keys, exists := data["keys"].([]interface{}); exists {

				for i, keyItem := range keys {
					if keyData, ok := keyItem.(map[string]interface{}); ok {
						wg.Add(1)
						go func(address, password, privateKey string) {
							defer wg.Done()
							ApiCheckBalance(address, password, privateKey, &addressWithTxs)
							mu.Lock()
							logger.Log(fmt.Sprintf("Total addresses: %d, Addresses with transactions: %d", totalAddresses-i, addressWithTxs), false)
							mu.Unlock()
						}(keyData["address"].(string), keyData["password"].(string), keyData["private_key"].(string))
						time.Sleep(20 * time.Millisecond)
					}
				}
			} else {
				logger.Log("Failed to read Addresses, please make sure you have scraped and a file called private_keys.json exists", true)
			}
		}
	}
	wg.Wait()
}

func ApiCheckBalance(address string, password string, privateKey string, addressWithTxs *int) {
	apiUrl := "https://api.blockchain.info/haskoin-store/btc/address/" + address + "/balance"

	resp, err := http.Get(apiUrl)
	if err != nil {
		logger.Log("Failed to get balance for address: "+address, true)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		time.Sleep(25 * time.Millisecond)
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Log("Failed to parse response body for address: "+address, true)
		return
	}
	//mu.Lock()
	if txs, ok := result["txs"].(float64); ok && txs > 0 {
		*addressWithTxs++
		result["password"] = password
		result["private_key"] = privateKey
		result["balance"] = satoshi_to_btc(result["received"].(float64))
		jsonPrettyWithKeys, _ := json.MarshalIndent(result, "", "  ")
		filemanager.WriteFile("balance_check.json", []string{string(jsonPrettyWithKeys)})
	}

	if result["confirmed"].(float64) > 0 {
		logger.Log(fmt.Sprintf("Address %s has a balance of %f BTC", address, result["confirmed"].(float64)), false)
	}
	//mu.Unlock()
}

func satoshi_to_btc(value float64) float64 {
	return value / 100000000.00000000
}
