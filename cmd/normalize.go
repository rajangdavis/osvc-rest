package cmd

import (
	"bytes"
	"fmt"
	"encoding/csv"
	"github.com/buger/jsonparser"
	"os"
	"strconv"
)

func iterateThroughRows(item []byte, arrayToMod [][]map[string]interface{}) [][]map[string]interface{} {
	val := item

	columnNames, _, _, _ := jsonparser.Get(val, "columnNames")
	rows, _, _, _ := jsonparser.Get(val, "rows")

	var itemArray []map[string]interface{}

	jsonparser.ArrayEach(rows, func(row []byte, dataType jsonparser.ValueType, offset int, err error) {

		resultsHash := make(map[string]interface{})

		columnIndex := 0
		jsonparser.ArrayEach(columnNames, func(column []byte, columnDataType jsonparser.ValueType, offset int, err error) {

			thisRow, _, _, err := jsonparser.Get(row, "["+strconv.Itoa(columnIndex)+"]")

			parsedColumn, _ := jsonparser.ParseString(column)
			if _, err := strconv.Atoi(string(thisRow)); err == nil {
				parsedRow, _ := jsonparser.ParseInt(thisRow)
				resultsHash[parsedColumn] = parsedRow
			} else {
				parsedRow, _ := jsonparser.ParseString(thisRow)
				resultsHash[parsedColumn] = parsedRow
			}
			var newIndex = columnIndex + 1
			columnIndex = newIndex
		})

		itemArray = append(itemArray, resultsHash)
	})

	arrayToMod = append(arrayToMod, itemArray)
	return arrayToMod
}

func csvReport(byteData []byte, csvName string){
	file, _ := os.Create(csvName + ".csv")
	defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()
	
	var stringColumns  []string

	columnNames, _, _, _ := jsonparser.Get(byteData, "columnNames")
	
	jsonparser.ArrayEach(columnNames, func(column []byte, dataType jsonparser.ValueType, offset int, err error) {
		parsedColumn, _ := jsonparser.ParseString(column)
		stringColumns = append(stringColumns, parsedColumn)
	})
	writer.Write(stringColumns)
	
	rows, _, _, _ := jsonparser.Get(byteData, "rows")

	rowCount := 0
	
	jsonparser.ArrayEach(rows, func(row []byte, dataType jsonparser.ValueType, offset int, err error) {
		var stringRows []string
		jsonparser.ArrayEach(row, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			parsedRow, _ := jsonparser.ParseString(value)
			stringRows = append(stringRows, parsedRow)
		})
		rowCount = rowCount + 1
		// TODO 

		// Change to recursive function
		// Keep writing to file
		writer.Write(stringRows)
	})

}

func normalizeReport(byteData []byte, jsonString []byte, results *[]map[string]interface{})  []map[string]interface{} {

	columnNames, _, _, _ := jsonparser.Get(byteData, "columnNames")
	rows, _, _, _ := jsonparser.Get(byteData, "rows")
	
	var itemArray []map[string]interface{}

	jsonparser.ArrayEach(rows, func(row []byte, dataType jsonparser.ValueType, offset int, err error) {

		resultsHash := make(map[string]interface{})

		columnIndex := 0

		jsonparser.ArrayEach(columnNames, func(column []byte, columnDataType jsonparser.ValueType, offset int, err error) {

			thisRow, _, _, err := jsonparser.Get(row, "["+strconv.Itoa(columnIndex)+"]")
			parsedColumn, _ := jsonparser.ParseString(column)

			if _, err := strconv.Atoi(string(thisRow)); err == nil {
				parsedRow, _ := jsonparser.ParseInt(thisRow)
				resultsHash[parsedColumn] = parsedRow
			} else {
				parsedRow, _ := jsonparser.ParseString(thisRow)
				resultsHash[parsedColumn] = parsedRow
			}

			var newIndex = columnIndex + 1
			columnIndex = newIndex
		})

		itemArray = append(itemArray, resultsHash)
	})

	for i := 0; i < len(itemArray); i++ {
		*results = append(*results, itemArray[i])
	}

	if(len(itemArray) == 10000 && len(*results) < 1000000){
		offset, _, _, _ := jsonparser.Get(jsonString, "offset")
		intOffset, _ := strconv.ParseInt(string(offset), 10, 64)
		intOffset = intOffset + 10000
		jsonString, _ = jsonparser.Set(jsonString, []byte(strconv.FormatInt(intOffset, 10)), "offset")
		updatedJsonData := bytes.NewBuffer(jsonString)
		bodyBytes := connect("POST", "analyticsReportResults", updatedJsonData)
		return normalizeReport(bodyBytes, jsonString, results)
	}else{
		return *results
	}

}

func normalizeQuery(byteData []byte) [][]map[string]interface{} {

	items, dataType, offset, err := jsonparser.Get(byteData, "items")

	if err != nil {
		_, _ = dataType, offset
		parsedError, _ := jsonparser.ParseString(byteData)
		fmt.Fprintf(os.Stdout, "%s", parsedError)
		os.Exit(0)
	}

	var results [][]map[string]interface{}

	jsonparser.ArrayEach(items, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		results = iterateThroughRows(value, results)
	})

	return results

}
