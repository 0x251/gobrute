package scraper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"l4tt/gobrute/modules/filemanager"
	"l4tt/gobrute/modules/logger"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ripemd160"
)

func Scrape(target string, passwordlist string) {
	passwords, err := filemanager.ReadFile(passwordlist)
	if err != nil {
		logger.Log("Error reading password list: "+err.Error(), true)
		return
	}

	totalPasswords := len(passwords)
	data := map[string]interface{}{
		"count": totalPasswords,
		"keys":  make(map[string]interface{}, totalPasswords),
	}

	filemanager.DeleteFile("results/private_keys.json")

	keysMap := data["keys"].(map[string]interface{})

	startTime := time.Now()

	for i, password := range passwords {
		privateKey, err := from_passphrase(password)
		if err != nil {
			logger.Log("Error generating private key: "+err.Error(), true)
			continue
		}

		pubKey := elliptic.Marshal(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
		pubKeyHash := hash160(pubKey)
		address := base58CheckEncode(pubKeyHash, 0x00)

		keysMap[fmt.Sprintf("key_%d", i)] = map[string]interface{}{
			"private_key": fmt.Sprintf("%x", privateKey.D.Bytes()),
			"password":    password,
			"public_key":  fmt.Sprintf("%x", pubKey),
			"address":     address,
		}

		remainingPasswords := totalPasswords - (i + 1)

		if remainingPasswords%100 == 0 || remainingPasswords < 100000000 {
			elapsed := time.Since(startTime)
			processed := i + 1
			var ttc time.Duration
			if processed > 0 {
				averageTime := elapsed / time.Duration(processed)
				ttc = averageTime * time.Duration(remainingPasswords)
			} else {
				ttc = 0
			}
			logger.Status(fmt.Sprintf("Total: [\033[34m%d\033[0m], Remaining: [\033[34m%d\033[0m], TTC: [\033[35m%s\033[0m]", totalPasswords, remainingPasswords, ttc.Round(time.Second)))
		}
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	if totalPasswords > 10 {
		logger.Log("Data compressed to not spam the console", false)
	} else {
		logger.Log(string(jsonData), false)
	}

	logger.Log("Saved to [\033[34mresults/private_keys.json\033[0m]", false)
	filemanager.WriteFile("results/private_keys.json", []string{string(jsonData)})
}

func from_passphrase(passphrase string) (*ecdsa.PrivateKey, error) {
	hash := sha256.Sum256([]byte(passphrase))

	privateKey := new(ecdsa.PrivateKey)
	privateKey.D = new(big.Int).SetBytes(hash[:])
	privateKey.PublicKey.Curve = btcec.S256()

	if privateKey.D.Cmp(btcec.S256().Params().N) >= 0 || privateKey.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key derived from passphrase")
	}

	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privateKey.D.Bytes())

	return privateKey, nil
}

func hash160(data []byte) []byte {
	hash := sha256.Sum256(data)
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash[:])
	return ripemd160.Sum(nil)
}

func base58CheckEncode(data []byte, version byte) string {
	versionedPayload := append([]byte{version}, data...)
	firstSHA := sha256.Sum256(versionedPayload)
	secondSHA := sha256.Sum256(firstSHA[:])
	checksum := secondSHA[:4]
	finalPayload := append(versionedPayload, checksum...)
	return base58Encode(finalPayload)
}

func base58Encode(input []byte) string {
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	num := new(big.Int).SetBytes(input)
	var result []byte
	mod := new(big.Int)

	for num.Sign() > 0 {
		num.DivMod(num, big.NewInt(58), mod)
		result = append([]byte{alphabet[mod.Int64()]}, result...)
	}

	for _, b := range input {
		if b == 0x00 {
			result = append([]byte{alphabet[0]}, result...)
		} else {
			break
		}
	}

	return string(result)
}
