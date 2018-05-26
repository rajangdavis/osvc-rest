package cmd

import (
	"github.com/buger/jsonparser"
	"fmt"
	"os"
	"strconv"
	"encoding/json"
)

func iterateThroughRows(item []byte, arrayToMod [][]map[string]interface{}) [][]map[string]interface{} {
	val := item
	
	columnNames, _, _, _ := jsonparser.Get(val,"columnNames")
	rows, _, _, _ := jsonparser.Get(val,"rows")

	var itemArray []map[string]interface{}

	jsonparser.ArrayEach(rows, func(row []byte, dataType jsonparser.ValueType, offset int, err error) {

		resultsHash := make(map[string]interface{})

		columnIndex := 0
		jsonparser.ArrayEach(columnNames, func(column []byte, columnDataType jsonparser.ValueType, offset int, err error) {

			thisRow, _, _, err := jsonparser.Get(row,"[" + strconv.Itoa(columnIndex) + "]")

    		parsedColumn, _ := jsonparser.ParseString(column)
    		if _, err := strconv.Atoi(string(thisRow)); err == nil{
    			parsedRow, _ := jsonparser.ParseInt(thisRow)
    			resultsHash[parsedColumn] = parsedRow
			}else{
				parsedRow, _ := jsonparser.ParseString(thisRow)
    			resultsHash[parsedColumn] = parsedRow
			}
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

    var results [][]map[string]interface{}
	var itemArray []map[string]interface{}

	jsonparser.ArrayEach(rows, func(row []byte, dataType jsonparser.ValueType, offset int, err error) {

		resultsHash := make(map[string]interface{})

		columnIndex := 0

		jsonparser.ArrayEach(columnNames, func(column []byte, columnDataType jsonparser.ValueType, offset int, err error) {

			thisRow, _, _, err := jsonparser.Get(row,"[" + strconv.Itoa(columnIndex) + "]")
    		parsedColumn, _ := jsonparser.ParseString(column)

    		if _, err := strconv.Atoi(string(thisRow)); err == nil{
    			parsedRow, _ := jsonparser.ParseInt(thisRow)
    			resultsHash[parsedColumn] = parsedRow
			}else{
				parsedRow, _ := jsonparser.ParseString(thisRow)
    			resultsHash[parsedColumn] = parsedRow
			}

    		var newIndex = columnIndex + 1
    		columnIndex = newIndex
    	})

		itemArray = append(itemArray,resultsHash)
	})

	results = append(results,itemArray)
    jsonData, _ := json.MarshalIndent(results[0],"","  ")
	fmt.Fprintf(os.Stdout, "%s", jsonData)

}

func normalizeQuery(byteData []byte) [][]map[string]interface{}{

	items, dataType, offset, err := jsonparser.Get(byteData,"items")

    if err!= nil{
    	_ , _ = dataType, offset
    	parsedError, _ := jsonparser.ParseString(byteData)
		fmt.Fprintf(os.Stdout, "%s", parsedError)
    	os.Exit(0)
    }   

    var results [][]map[string]interface{}

    jsonparser.ArrayEach(items, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
    	results = iterateThroughRows(value,results)
	})

	return results    

}