package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"encoding/json"
	"net/url"
)

func pqueryCheck(args []string) []string {
	queriesArray := []string{}
	if len(args) <= 1{
		fmt.Println("Error: Must set at least two queries")
		os.Exit(0)
	}else{
		for i := 0; i < len(args); i++ {
			queriesArray = append(queriesArray,args[i])
		}
	}
	return queriesArray
}

func asyncQuery(resourceUrl string, ch chan <- []map[string]interface{}){
	bodyBytes := connect("GET", resourceUrl, nil)
	initResult := normalizeQuery(bodyBytes)
	ch <- initResult[0]
}

func runParallelQueries(urls []string) <- chan []map[string]interface{}{

	ch := make(chan []map[string]interface{}, len(urls))

	for i := 0; i < len(urls); i++ {

		queryUrl := "queryResults?query=" + url.PathEscape(urls[i])
		go asyncQuery(queryUrl,ch)
	}

	return ch
} 

func printParallelQueries(cmd *cobra.Command, args []string) error {
	
	var queriesToRun = pqueryCheck(args)
	results := runParallelQueries(queriesToRun)
	var finalResults = make([]map[string]interface{}, 0)

	for _ = range queriesToRun {
		result := <-results
		for i := 0; i < len(result); i++ {
			finalResults = append(finalResults,result[i])
		}
	}

	jsonData, _ := json.MarshalIndent(finalResults,"","  ")
	fmt.Fprintf(os.Stdout, "%s", jsonData)
	return nil
	
}

var pquery = &cobra.Command{
	Use: "pquery",
	Short: "Runs multiple ROQL queries in parallel",
	Long: "\033[93mRuns one or more ROQL queries and returns parsed results\033[0m \033[0;32m\n\nSingle Query Example: \033[0m \n$ osvc-rest query \"DESCRIBE\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE \033[0;32m\n\nMultiple Queries Example:\033[0m \n$ osvc-rest query \"SELECT * FROM INCIDENTS LIMIT 100\" \"SELECT * FROM SERVICEPRODUCTS LIMIT 100\" \"SELECT * FROM SERVICECATEGORIES LIMIT 100\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE",
	RunE: printParallelQueries,
}

func init(){
	RootCmd.AddCommand(pquery)
}