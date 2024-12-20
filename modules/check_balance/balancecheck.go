package check_balance

import (
	"encoding/json"
	"fmt"
	"l4tt/gobrute/modules/filemanager"
	"l4tt/gobrute/modules/logger"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	apiURLTemplate = "https://api.blockchain.info/haskoin-store/btc/address/%s/balance"
)

func CheckBalance(concurrent_threads int) {

	if concurrent_threads <= 0 {
		concurrent_threads = 10
	}

	if concurrent_threads > 10 {
		logger.Log("Warning: Using more than 10 concurrent threads may cause rate limiting or connection errors.", true)
	}

	if concurrent_threads > 150 {
		logger.Log("Too many concurrent threads, please use a lower number", true)
		return
	}

	jsonData, err := filemanager.ReadJSON("results/private_keys.json")
	if err != nil {
		logger.Log("Failed to read Addresses. Ensure 'private_keys.json' exists after scraping.", true)
		return
	}

	filemanager.DeleteFile("results/balance_check.json")

	var wg sync.WaitGroup
	var addressWithTxs int32
	var addressBalance int32
	var errorCount int32
	var processedAddresses int32

	startTime := time.Now()

	for _, item := range jsonData {
		data, ok := item.(map[string]interface{})
		if !ok {
			logger.Log("Invalid data format in JSON.", true)
			continue
		}

		totalAddresses, ok := data["count"].(float64)
		if !ok {
			logger.Log("Invalid 'count' value in JSON.", true)
			continue
		}

		keysMap, exists := data["keys"].(map[string]interface{})
		if !exists {
			logger.Log("Missing 'keys' in JSON data.", true)
			continue
		}

		sem := make(chan struct{}, concurrent_threads)

		index := 0
		for _, keyItem := range keysMap {
			keyData, ok := keyItem.(map[string]interface{})
			if !ok {
				logger.Log("Invalid key item format.", true)
				continue
			}

			address, aOk := keyData["address"].(string)
			password, pOk := keyData["password"].(string)
			privateKey, kOk := keyData["private_key"].(string)

			if !aOk || !pOk || !kOk {
				logger.Log("Incomplete key data received.", true)
				continue
			}

			wg.Add(1)
			sem <- struct{}{}
			go func(address, password, privateKey string, index int) {
				defer wg.Done()
				defer func() { <-sem }()

				ApiCheckBalance(address, password, privateKey, &addressWithTxs, &addressBalance, &errorCount)

				currentProcessed := atomic.AddInt32(&processedAddresses, 1)
				addressesLeft := int(totalAddresses) - int(currentProcessed)
				elapsed := time.Since(startTime)
				var ttc time.Duration
				if currentProcessed > 0 {
					averageTimePerAddress := elapsed / time.Duration(currentProcessed)
					ttc = averageTimePerAddress * time.Duration(addressesLeft)
				} else {
					ttc = 0
				}

				currentWithTxs := atomic.LoadInt32(&addressWithTxs)
				currentBalance := atomic.LoadInt32(&addressBalance)
				currentError := atomic.LoadInt32(&errorCount)

				logger.Status(fmt.Sprintf(
					"Total: [\033[34m%d\033[0m], With Txs: [\033[33m%d\033[0m], With Balance: [\033[32m%d\033[0m], CErrors: [\033[31m%d\033[0m] TTC: [\033[35m %s \033[0m]",
					int(totalAddresses)-int(currentProcessed), currentWithTxs, currentBalance, currentError, ttc.Round(time.Second)))
			}(address, password, privateKey, index)

			index++
		}
	}

	wg.Wait()
	logger.Log("Balance check completed successfully results in [\033[34mresults/balance_check.json\033[0m]", false)
}

func ApiCheckBalance(address, password, privateKey string, addressWithTxs *int32, addressBalance *int32, errorCount *int32) {
	apiURL := fmt.Sprintf(apiURLTemplate, address)

	var resp *http.Response
	var err error

	for {
		resp, err = http.Get(apiURL)
		if err != nil {
			atomic.AddInt32(errorCount, 1)
			return
		}

		if resp.StatusCode == http.StatusOK {
			break
		} else {
			atomic.AddInt32(errorCount, 1)
			resp.Body.Close()
			time.Sleep(25 * time.Millisecond)
		}
	}
	defer resp.Body.Close()

	var result struct {
		Txs       float64 `json:"txs"`
		Received  float64 `json:"received"`
		Confirmed float64 `json:"confirmed"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Log(fmt.Sprintf("Failed to parse response for address: %s", address), true)
		atomic.AddInt32(errorCount, 1)
		return
	}

	if result.Received != 5460 && result.Txs > 0 {
		atomic.AddInt32(addressWithTxs, 1)

		btcReceived := satoshiToBTC(result.Received)
		btcConfirmed := satoshiToBTC(result.Confirmed)

		if result.Confirmed > 0 {
			atomic.AddInt32(addressBalance, 1)
			logger.Log(fmt.Sprintf("Address: %s, Balance: %f", address, btcConfirmed), false)
		}

		record := map[string]interface{}{
			"address":            address,
			"txs":                result.Txs,
			"received":           result.Received,
			"confirmed":          result.Confirmed,
			"formated_received":  btcReceived,
			"formated_confirmed": btcConfirmed,
			"password":           password,
			"private_key":        privateKey,
		}

		jsonData, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			logger.Log(fmt.Sprintf("Failed to marshal JSON for address: %s", address), true)
			return
		}

		if err := filemanager.AppendToFile("results/balance_check.json", string(jsonData)+"\n"); err != nil {
			logger.Log(fmt.Sprintf("Failed to write to balance_check.json for address: %s", address), true)
		}

	}
}

func satoshiToBTC(value float64) float64 {
	return value / 1e8
}
