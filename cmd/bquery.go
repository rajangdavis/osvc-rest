package cmd

import (
	"encoding/json"
	"strconv"
	"strings"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
)

var queryToCount, bulkTableName string
var batchCount, numQueries int

func checkBulkQueryFlags(flags *pflag.FlagSet) error {
	interfaceAndPassword()
	if queryToCount == "" || bulkTableName == "" {
		fmt.Println("\033[31mError: Must use both the --select and --from flag with for performing bulk queries\033[0m ")
		os.Exit(0)
	}

	return nil
}

func splitQuery(queryToSplit string, bulkTableName string) []string{
	splitQuery := strings.Split(queryToCount, ",")
	
	var queriesToCount  []string
	
	for i := 0; i < len(splitQuery); i++ {
		splitAs := strings.Split(splitQuery[i], " as ")
		countQuery := "SELECT DISTINCT COUNT(" + splitAs[0] + ") as count FROM " + bulkTableName
		queriesToCount = append(queriesToCount, countQuery)
	}

	return queriesToCount
}


func runBulkQuery(cmd *cobra.Command, args []string) error {
	
	var queriesToCount = splitQuery(queryToCount, bulkTableName)
	results := runParallelQueries(queriesToCount)
	var maxValue int64

	for _ = range queriesToCount {
		result := <-results
		for i := 0; i < len(result); i++ {
			innerCount := result[i][0]["count"]
			int64InnerCount := innerCount.(int64)
			if(maxValue < int64InnerCount){
				maxValue = int64InnerCount
			}
		}
	}

	fmt.Fprintf(os.Stderr, "%s", "\nFetching " + strconv.Itoa(int(maxValue)) + " rows\n\n")

	leftOver := maxValue % int64(batchCount)
	var numberOfRequests int64
	if leftOver > 0{
		numberOfRequests = ((maxValue - leftOver)/int64(batchCount) + 1)
	}else{
		numberOfRequests = (maxValue - leftOver)/int64(batchCount)
	}

	var queriesToRun []string

	for i := 0; i < int(numberOfRequests); i++ {
		queryToFetch := "SELECT " + queryToCount + " FROM " + bulkTableName + " LIMIT " + strconv.Itoa(batchCount) + " OFFSET " + strconv.Itoa(batchCount * i)
		queriesToRun = append(queriesToRun, queryToFetch)
	}

	var lowerBound int

	remainderQueries := len(queriesToRun) % numQueries

	numberOfBatchedQueries := ((len(queriesToRun) - remainderQueries)/numQueries)

	upperBound := numberOfBatchedQueries * numQueries

	var finalResults = make([]map[string]interface{}, 0)
	
	for i := 0; i < numberOfBatchedQueries; i++ {

		innerUpperBound := lowerBound + numQueries

		var currentQuerySet []string

		if(innerUpperBound > upperBound){
			currentQuerySet = queriesToRun[lowerBound:]
		}else{
			currentQuerySet = queriesToRun[lowerBound:innerUpperBound]
		}

		results := runParallelQueries(currentQuerySet)

		for _ = range currentQuerySet {
			result := <-results
			for i := 0; i < len(result[0]); i++ {
				finalResults = append(finalResults, result[0][i])
			}
		}

		lowerBound = lowerBound + (numQueries)
	}
	jsonData, _ := json.MarshalIndent(finalResults, "", "  ")
	fmt.Fprintf(os.Stdout, "%s", jsonData)

	return nil
}

var bquery = &cobra.Command{
	Use:   "bquery",
	Short: "Uses concurrency and ROQL behind the scenes to bulk fetch data",
	// Long:  "\033[93mRuns one or more ROQL queries and returns parsed results\033[0m \033[0;32m\n\nSingle Query Example: \033[0m \n$ osvc-rest query \"DESCRIBE\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE \033[0;32m\n\nMultiple Queries Example:\033[0m \n$ osvc-rest query \"SELECT * FROM INCIDENTS LIMIT 100\" \"SELECT * FROM SERVICEPRODUCTS LIMIT 100\" \"SELECT * FROM SERVICECATEGORIES LIMIT 100\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE \033[0;32m\n\nParallel Queries Example:\033[0m \n$ osvc-rest query --parallel \"SELECT * FROM INCIDENTS LIMIT 20000\" \"SELECT * FROM INCIDENTS Limit 20000 OFFSET 20000\" \"SELECT * FROM INCIDENTS Limit 20000 OFFSET 40000\" \"SELECT * FROM INCIDENTS Limit 20000 OFFSET 60000\" \"SELECT * FROM INCIDENTS Limit 20000 OFFSET 80000\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkBulkQueryFlags(cmd.Flags())
	},
	RunE:  runBulkQuery,
}

func init() {
	bquery.Flags().StringVarP(&queryToCount, "select", "", "", "The main query that is getting bulk requests performed on")
	bquery.Flags().StringVarP(&bulkTableName, "from", "", "", "The table that is getting queried against")
	bquery.Flags().IntVarP(&batchCount, "batch", "", 10000, "How many rows to batch per request")
	bquery.Flags().IntVarP(&numQueries, "group", "", 4, "How many queries to run in parallel")
	RootCmd.AddCommand(bquery)
}
