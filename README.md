# gobrute
Bitcoin brain wallet brute forcer written in Golang

### Overview
gobrute is a high-performance Bitcoin brain wallet brute-forcing tool written in Go. It allows you to target and brute-force brain wallet addresses efficiently while providing options to scrape and verify wallet balances at an impressive speed.

# Features
- [x] --target scrape: Scrape brain wallet candidates.
- [x] --target checkbalance: Verify wallet balances.
- [x] --passwordlist: Password list
- [x] Speed: Achieves up to 1,000 CPM.

Commands
```

gobrute --target scrape --passwordlist password_list.txt (Will make a private_keys.json, with the password & address & private key) 
```
# To check wallet balances
```
gobrute --target checkbalance
```
# Installation
To build and run gobrute, ensure you have Go installed. Then follow these steps:

```
cd gobrute

go build .

./gobrute --target [scrape|checkbalance]
```

### Disclaimer
This tool is provided for educational and research purposes only. Using it to target unauthorized wallets or engage in malicious activity is illegal and unethical.
