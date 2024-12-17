package scraper

import (
	"fmt"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"math/big"
	"l4tt/gobrute/modules/logger"
	"l4tt/gobrute/modules/filemanager"
	"encoding/json"
	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcd/btcec"
)
func Scrape(target string, passwordlist string) {
	passwords, err := filemanager.ReadFile(passwordlist)
	if err != nil {
		logger.Log("Error reading password list: "+err.Error(), true)
		return
	}

	data := map[string]interface{}{
		"count": len(passwords),
		"keys":  []map[string]interface{}{},
	}

	filemanager.DeleteFile("private_keys.json")

	for _, password := range passwords {
		privateKey, err := from_passphrase(password)
		
		if err != nil {
			logger.Log("Error generating private key: "+err.Error(), true)
			continue
		}
		pubKey := elliptic.Marshal(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
		pubKeyHash := hash160(pubKey)
		address := base58CheckEncode(pubKeyHash, 0x00)

		data["count"] = len(passwords)
		data["keys"] = append(data["keys"].([]map[string]interface{}), map[string]interface{}{
			"private_key": fmt.Sprintf("%x", privateKey.D.Bytes()),
			"password":    password,
			"public_key":  fmt.Sprintf("%x", pubKey),
			"address":     address,
		})
	}

	jsonData, _ := json.MarshalIndent(data, "", "  ")
	if len(passwords) > 10 {
		logger.Log("Data compressed to not spam the console", false)
	} else {
		logger.Log(string(jsonData), false)
	}

	logger.Log("Saved to private_keys.json", false)
	filemanager.WriteFile("private_keys.json", []string{string(jsonData)})
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


