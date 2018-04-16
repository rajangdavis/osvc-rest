package cmd

import (
	"fmt"
	"strings"
	"errors"
	"net/url"
	"os"
	"strconv"
	"io/ioutil"
	"net/http"
	"crypto/tls"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/buger/jsonparser"
)


func runQuery(cmd *cobra.Command, args []string) error {

	queryInit := []string{}

	queryFinal := ""

	if len(args) == 0{
		fmt.Println("Error: Must set at least one query")
		os.Exit(0)
	}else if len(args) == 1{
		queryFinal = url.PathEscape(args[0])
	}else{
		for i := 0; i < len(args); i++ {
			queryInit = append(queryInit,args[i])
		}
		queryFinal = url.PathEscape(strings.Join(queryInit, ";"))
	}


	queryUrl := "queryResults?query=" + queryFinal






	// put into connect.go

	if annotation != "" && len(annotation) > 40 {
		fmt.Println("Error: Annotation cannot be greater than 40 characters.")
		os.Exit(0)
	}else if version == "v1.4" && annotation == ""{
		fmt.Println("Error: An Annotation must be set when using CCOM version v1.4 (e.g. -a \"40 character annotation\")")
		os.Exit(0)
	}

	domain :=""

	if demoSite == true{
		domain = "rightnowdemo"
	}else{
		domain = "custhelp"
	}

	var client = &http.Client{}
    
    if noSslVerify == true{
	    tr := &http.Transport{
	        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	    }
	    client = &http.Client{Transport: tr}
    }

	var url = "https://"+ interfaceName +"." + domain +".com/services/rest/connect/" + version + "/"
	var finalUrl = url + queryUrl
  	req, err := http.NewRequest("GET", finalUrl, nil)
    req.Header.Add("Authorization","Basic " + basicAuth(userName,password))

    if version == "v1.4" && annotation != "" && len(annotation) <= 40{
	    req.Header.Add("OSvC-CREST-Application-Context", annotation)	
    }

    if utcTime == true{
	    req.Header.Add("OSvC-CREST-Time-UTC", "true")
    }

    rs, err := client.Do(req)
    // Process response
    if err != nil {
        fmt.Println("Error: Could not connect to site '" + url + "'")
        os.Exit(1)
    }
    defer rs.Body.Close()
 
    bodyBytes, err := ioutil.ReadAll(rs.Body)
    
    if err != nil {
        panic(err)
    }

    // return bodyBytes













 	// put into normalize.go
 	// handle JSON
    items, dataType, offset, err := jsonparser.Get(bodyBytes,"items")

    if err!= nil{
    	_ , _ = dataType, offset
    	parsedError, _ := jsonparser.ParseString(bodyBytes)
		fmt.Fprintf(os.Stdout, "%s", parsedError)
    	os.Exit(0)
    }   

    var results [][]map[string]string

    jsonparser.ArrayEach(items, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
    	val := value
		columnNames, dataType, offset, err := jsonparser.Get(val,"columnNames")
		rows, dataType, offset, err := jsonparser.Get(val,"rows")

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
    	
	})

    if len(results) == 1{
    	jsonData, _ := json.MarshalIndent(results[0],"","  ")
		fmt.Fprintf(os.Stdout, "%s", jsonData)
	}else{
    	jsonData, _ := json.MarshalIndent(results,"","  ")
		fmt.Fprintf(os.Stdout, "%s", jsonData)
	}


	return nil
}


func CheckRequiredFlags(flags *pflag.FlagSet) error {
	requiredError := false
	flagName := ""

	flags.VisitAll(func(flag *pflag.Flag) {
		requiredAnnotation := flag.Annotations[cobra.BashCompOneRequiredFlag]
		if len(requiredAnnotation) == 0 {
			return
		}

		flagRequired := requiredAnnotation[0] == "true"

		if flagRequired && !flag.Changed {
			requiredError = true
			flagName = flag.Name
		}
	})

	if requiredError {
		return errors.New("Required flag `" + flagName + "` has not been set")
	}

	return nil
}

// query represents the query command
var query = &cobra.Command{
	Use: "query",
	Short: "Runs one or more ROQL queries",
	Long: "Runs one or more ROQL queries",
	PreRunE: func(cmd *cobra.Command, args []string) error {		
		return CheckRequiredFlags(cmd.Flags())
	},
	RunE: runQuery,
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&userName,"username","u","", "Username to use for basic authentication")
    RootCmd.MarkPersistentFlagRequired("username")
    RootCmd.PersistentFlags().StringVarP(&password,"password","p","", "Password to use for basic authentication")
    RootCmd.MarkPersistentFlagRequired("password")
    RootCmd.PersistentFlags().StringVarP(&interfaceName,"interface","i","", "Oracle Service Cloud Interface to connect with")
    RootCmd.MarkPersistentFlagRequired("interface")
    
    RootCmd.PersistentFlags().BoolVarP(&demoSite,"demosite","",false, "Change the domain from 'custhelp' to 'rightnowdemo'")
    RootCmd.PersistentFlags().BoolVarP(&suppressRules,"suppress-rules","s",false, "Adds a header to suppress business rules")
    RootCmd.PersistentFlags().BoolVarP(&noSslVerify,"no-ssl-verify","",false, "Turns off SSL verification")
    RootCmd.PersistentFlags().StringVarP(&version,"version","v","v1.3", "Changes the CCOM version")
    RootCmd.PersistentFlags().StringVarP(&annotation,"annotate","a","", "Adds a custom header that adds an annotation")
    RootCmd.PersistentFlags().BoolVarP(&excludeNull,"exclude-null","e",false, "Adds a custom header to excludes null from results")
    RootCmd.PersistentFlags().BoolVarP(&utcTime,"utcTime","t",false, "Adds a custom header to return results using Coordinated Universal Time (UTC) format for time")

	RootCmd.AddCommand(query)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// query.PersistentFlags().String("foo", "", "A help for foo")
}	