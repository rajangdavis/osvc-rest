package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

func pqueryCheck(args []string) []string {
	queriesArray := []string{}
	if len(args) < 2 {
		fmt.Println("\033[31mError: Must set at least two queries \033[0m ")
		os.Exit(0)
	}
	for i := 0; i < len(args); i++ {
		queriesArray = append(queriesArray, args[i])
	}
	return queriesArray
}

func asyncQuery(resourceUrl string, ch chan<- [][]map[string]interface{}) {
	bodyBytes := connect("GET", resourceUrl, nil)
	initResult := normalizeQuery(bodyBytes)
	ch <- initResult
}

func runParallelQueries(urls []string) <-chan [][]map[string]interface{} {

	ch := make(chan [][]map[string]interface{}, len(urls))

	for i := 0; i < len(urls); i++ {

		queryUrl := "queryResults?query=" + url.PathEscape(urls[i])
		go asyncQuery(queryUrl, ch)
	}

	return ch
}

func printParallelQueries(args []string) error {

	var queriesToRun = pqueryCheck(args)
	results := runParallelQueries(queriesToRun)
	var finalResults = make([][]map[string]interface{}, 0)

	for _ = range queriesToRun {
		result := <-results
		for i := 0; i < len(result); i++ {
			finalResults = append(finalResults, result[i])
		}
	}

	jsonData, _ := json.MarshalIndent(finalResults, "", "  ")
	fmt.Fprintf(os.Stdout, "%s", jsonData)
	return nil

}
