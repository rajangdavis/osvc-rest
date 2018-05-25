package cmd

import (
	"fmt"
	"strings"
	"os"
	"io"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"encoding/json"
)

func httpCheck(args []string) []string {
	resourceUrls := []string{}
	if len(args) == 0{
		fmt.Println("Error: Must set at least one resource url")
		os.Exit(0)
	}else{
		for i := 0; i < len(args); i++ {
			resourceUrls = append(resourceUrls,args[i])
		}
	}
	return resourceUrls
}

func makeRequest(verb string, url string, optionalJson io.Reader, ch chan <-[]byte) {
	byteData := connect(verb,url,optionalJson)
	m, ok := gjson.Parse(string(byteData)).Value().(map[string]interface{})
	if !ok {
        fmt.Println("Error")
    }

	jsonData, _ := json.MarshalIndent(m,"","  ")

	ch <- jsonData
}

func runHttp(cmd *cobra.Command, args []string) error {

	resourceUrls := httpCheck(args)
	resourceUrlsCount := len(resourceUrls) 
	httpVerb := strings.ToUpper(cmd.Use)

	ch := make(chan []byte)
	
	for i := 0; i < resourceUrlsCount; i++ {
		go makeRequest(httpVerb,resourceUrls[i], nil, ch)	

		if resourceUrlsCount > 1{
			fmt.Fprintf(os.Stdout, "\n")
		}
		
		fmt.Fprintf(os.Stdout, "%s", <-ch)
	}

	return nil
}

var get = &cobra.Command{
	Use: "get",
	Short: "Runs one or more GET requests",
	Long: "\033[93mRuns one or more GET requests and returns parsed results\033[0m \033[0;32m\n\nSingle Query Example: \033[0m \n$ osvc-rest query \"DESCRIBE\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE \033[0;32m\n\nMultiple Queries Example:\033[0m \n$ osvc-rest query \"SELECT * FROM INCIDENTS LIMIT 100\" \"SELECT * FROM SERVICEPRODUCTS LIMIT 100\" \"SELECT * FROM SERVICECATEGORIES LIMIT 100\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE",
	RunE: runHttp,
}


// TODO 
	// update 
		// http get
		
	// add functionality for 
		// post
		// puts
		// delete
		
		// file attachments
		// exclude null data

	// add validations

func init(){
	RootCmd.AddCommand(get)
}