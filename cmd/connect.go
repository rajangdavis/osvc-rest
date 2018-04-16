package cmd

import (
	"fmt"
	"os"
	"net/http"
	"crypto/tls"
	"io/ioutil"

	// "strings"
	// "errors"
	// "net/url"
	// "strconv"
	// "encoding/json"
	// "github.com/spf13/cobra"
	// "github.com/spf13/pflag"
	// "github.com/buger/jsonparser"
)

func connect(requestType string,requestUrl string) []byte{

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
	var finalUrl = url + requestUrl
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

    return bodyBytes

}