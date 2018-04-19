package cmd

import (
	"github.com/buger/jsonparser"
	"fmt"
	"os"
	"strconv"
	"encoding/json"
)

func iterateThroughRows(item []byte, arrayToMod [][]map[string]string) [][]map[string]string{
	val := item
	
	columnNames, _, _, _ := jsonparser.Get(val,"columnNames")
	rows, _, _, _ := jsonparser.Get(val,"rows")

	var itemArray []map[string]string

	jsonparser.ArrayEach(rows, func(row []byte, dataType jsonparser.ValueType, offset int, err error) {

		resultsHash := make(map[string]string)

		columnIndex := 0
		jsonparser.ArrayEach(columnNames, func(column []byte, columnDataType jsonparser.ValueType, offset int, err error) {

			thisRow, _, _, err := jsonparser.Get(row,"[" + strconv.Itoa(columnIndex) + "]")
    		parsedColumn, _ := jsonparser.ParseString(column)
    		parsedRow, _ := jsonparser.ParseString(thisRow)
    		resultsHash[parsedColumn] = parsedRow
    		var newIndex = columnIndex + 1
    		columnIndex = newIndex
    	})

		itemArray = append(itemArray,resultsHash)
	})

	arrayToMod = append(arrayToMod,itemArray)
	return arrayToMod
}

func normalizeReport(byteData []byte){

    columnNames, _, _, _ := jsonparser.Get(byteData,"columnNames")
	rows, _, _, _ := jsonparser.Get(byteData,"rows")

    var results [][]map[string]string
	var itemArray []map[string]string

	jsonparser.ArrayEach(rows, func(row []byte, dataType jsonparser.ValueType, offset int, err error) {

		resultsHash := make(map[string]string)

		columnIndex := 0
		jsonparser.ArrayEach(columnNames, func(column []byte, columnDataType jsonparser.ValueType, offset int, err error) {

			thisRow, _, _, err := jsonparser.Get(row,"[" + strconv.Itoa(columnIndex) + "]")
    		parsedColumn, _ := jsonparser.ParseString(column)
    		parsedRow, _ := jsonparser.ParseString(thisRow)
    		resultsHash[parsedColumn] = parsedRow
    		var newIndex = columnIndex + 1
    		columnIndex = newIndex
    	})

		itemArray = append(itemArray,resultsHash)
	})

	results = append(results,itemArray)
    jsonData, _ := json.MarshalIndent(results[0],"","  ")
	fmt.Fprintf(os.Stdout, "%s", jsonData)

}

func normalizeQuery(byteData []byte){

	items, dataType, offset, err := jsonparser.Get(byteData,"items")

    if err!= nil{
    	_ , _ = dataType, offset
    	parsedError, _ := jsonparser.ParseString(byteData)
		fmt.Fprintf(os.Stdout, "%s", parsedError)
    	os.Exit(0)
    }   

    var results [][]map[string]string

    jsonparser.ArrayEach(items, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
    	results = iterateThroughRows(value,results)
	})

    if len(results) == 1{
    	jsonData, _ := json.MarshalIndent(results[0],"","  ")
		fmt.Fprintf(os.Stdout, "%s", jsonData)
	}else{
    	jsonData, _ := json.MarshalIndent(results,"","  ")
		fmt.Fprintf(os.Stdout, "%s", jsonData)
	}

}