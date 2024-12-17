package filemanager

import (
	"bufio"
	"encoding/json"
	"os"
	"l4tt/gobrute/modules/logger"

)

func ReadFile(filename string) ([]string, error) {
	logger.Log("Reading file: "+filename, false)
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
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var data []interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}


func DeleteFile(filename string) error {
	return os.Remove(filename)
}
