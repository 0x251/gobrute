package filemanager

import (
	"bufio"
	"encoding/json"
	"l4tt/gobrute/modules/logger"
	"os"
)

func ReadFile(filename string) ([]string, error) {
	logger.Log("Reading file: [\033[34m"+filename+"\033[0m]", false)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func WriteFile(filename string, lines []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data := "["
	for i, line := range lines {
		if i > 0 {
			data += ","
		}
		data += line
	}
	data += "]"

	_, err = file.WriteString(data + "\n")
	return err
}

func ReadJSON(filename string) ([]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(bufio.NewReader(file))
	var data []interface{}
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func DeleteFile(filename string) error {
	return os.Remove(filename)
}

func AppendToFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(data)
	return err
}

func WritePassword(filename string, data []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range data {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
