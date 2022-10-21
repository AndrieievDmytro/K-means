package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func readCsv(path string) ([][]string, error) {
	var dataFile *os.File
	var err error

	dataFile, err = os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer dataFile.Close()
	if err == nil {
		var buf []byte
		var rerr error
		// Reading text from file
		buf, rerr = ioutil.ReadAll(dataFile)
		if rerr != nil {
			return nil, err
		}
		// Parsing from comma-separated text
		r := csv.NewReader(strings.NewReader(string(buf)))
		records, err := r.ReadAll()
		if err != nil {
			return nil, err
		}
		return records, nil
	}
	return nil, err
}

func convertStrArrayToJson(records [][]string) string {
	// Converting from array of string to JSON
	jsonData := ""
	strNum := 0
	paramsLength := len(records[0]) - 1
	for _, record := range records {
		wrongStr := false
		if len(record) < paramsLength || len(record) > paramsLength+1 {
			fmt.Println("Wrong parameters count")
			wrongStr = true
		}
		flName := record[len(record)-1] // Cutting flower name
		record = record[:len(record)-1]
		strNum++
		stringArray := "["                // Opening sq bracket
		for _, arrField := range record { // Filling string representation of array
			_, err := strconv.ParseFloat(arrField, 64)
			if err != nil {
				fmt.Println("Wrong parameters type in string: ", strNum)
				wrongStr = true
			}
			stringArray += arrField + ","
		}
		stringArray = stringArray[:len(stringArray)-1] // Removing last ',' character
		stringArray += "]"                             // Closing sq bracket
		if !wrongStr {
			jsonData += "{ \"name\": \"" + flName + "\", \"params\":" + stringArray + ", \"Distance\": [] }," // Converting to JSON
		}
	}
	jsonData = "[" + jsonData[:len(jsonData)-1] + "]"
	return jsonData

}

func (f *Flowers) readData(path string) {
	records, err := readCsv(path)
	if err == nil {
		jsonData := convertStrArrayToJson(records)
		json.Unmarshal([]byte(jsonData), &f.Fl)
	} else {
		fmt.Println("Read CSV error: " + err.Error())
		os.Exit(1)
	}
}
