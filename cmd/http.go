package cmd

import (
	"fmt"
	"strings"
	"os"
	"io"
	"net/url"
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
			urlToEncode := url.PathEscape(strings.Replace(args[i],"?","XXX_REPLACE_QUESTION_MARK_XXX",-1))
			resourceUrls = append(resourceUrls, strings.Replace(urlToEncode,"XXX_REPLACE_QUESTION_MARK_XXX","?",-1))
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
	Short: "Performs one or more GET requests",
	Long: "\033[93mPerforms one or more GET requests and returns parsed results\033[0m \033[0;32m\n\nSingle Query Example: \033[0m \n$ osvc-rest query \"DESCRIBE\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE \033[0;32m\n\nMultiple Queries Example:\033[0m \n$ osvc-rest query \"SELECT * FROM INCIDENTS LIMIT 100\" \"SELECT * FROM SERVICEPRODUCTS LIMIT 100\" \"SELECT * FROM SERVICECATEGORIES LIMIT 100\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE",
	RunE: runHttp,
}




var lookupName, filters string


func checkPostPutFlags(flags *pflag.FlagSet) error {

	if data == "" && id == 0 {
		fmt.Println("\033[31mError: Must use either the --name or --id flag with for working with the AnalyticsReportResults object")
		os.Exit(0)
	}

	return nil
}

var post = &cobra.Command{
	Use: "post",
	Short: "Performs a POST request",
	Long: "\033[93mPerforms a POST request and returns parsed results\033[0m \033[0;32m\n\nExample: \033[0m \n$ osvc-rest post \"incidents\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE\n\n",
	PreRunE:func(cmd *cobra.Command, args []string) error {		
		return checkPostPutFlags(cmd.Flags())
	},
	RunE: runHttp,
}

func init(){
	RootCmd.AddCommand(get)
	post.Flags().StringVarP(&filters,"filters","f","", "Adds filters for reporting")
	post.Flags().StringVarP(&lookupName,"name","n","", "Sets the lookupName of the AnalyticsReport that we wish to run")
	RootCmd.AddCommand(post)
}