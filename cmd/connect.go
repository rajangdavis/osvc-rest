package cmd

import (
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
)

func buildRequest(method, requestUrl string) (*http.Request, error, string){
	var url = "https://"+ interfaceName +"." + setDomain() +".com/services/rest/connect/v" + version + "/"
	var finalUrl = url + requestUrl
  	req, err := http.NewRequest("GET", finalUrl, nil)
    req.Header.Add("Authorization","Basic " + basicAuth(userName,password))

    if version == "1.4" && annotation != "" && len(annotation) <= 40{
	    req.Header.Add("OSvC-CREST-Application-Context", annotation)	
    }

    if utcTime == true{
	    req.Header.Add("OSvC-CREST-Time-UTC", "true")
    }
    return req, err, url
}

func connect(requestType string,requestUrl string) []byte{
	checkAnnotation()
	var client = checkSSL()
	req,err, url := buildRequest(requestType,requestUrl)
    rs, err := client.Do(req)

    if err != nil {
        fmt.Println("\033[31mError: Could not connect to site '" + url + "'")
        os.Exit(1)
    }
    defer rs.Body.Close()
 
    bodyBytes, err := ioutil.ReadAll(rs.Body)
    
    if err != nil {
        panic(err)
    }

    return bodyBytes
}