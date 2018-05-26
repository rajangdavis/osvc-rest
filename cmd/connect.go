package cmd

import (
    "fmt"
	"os"
	"net/http"
    "io"
	"io/ioutil"
    "strings"
)

func buildRequest(method string, requestUrl string, jsonData io.Reader) (*http.Request, error, string){
	var url = "https://"+ interfaceName +"." + setDomain() +".com/services/rest/connect/" + version + "/"
	var finalUrl = url + strings.Replace(requestUrl," ","%20",-1)

    var req *http.Request
    var err error

    if method == "POST" || method == "PATCH"{
        req, err = http.NewRequest("POST", finalUrl, jsonData)
        req.Header.Set("Content-Type", "application/json")
        if method == "PATCH"{
            req.Header.Set("X-HTTP-Method-Override", "PATCH")
        }
    }else{
        req, err = http.NewRequest(method, finalUrl, nil)
    }

    req.Header.Add("Authorization","Basic " + basicAuth(userName,password))

    if (version == "v1.4" || version == "latest") && annotation != "" && len(annotation) <= 40{
	    req.Header.Add("OSvC-CREST-Application-Context", annotation)	
    }

    if utcTime == true{
	    req.Header.Add("OSvC-CREST-Time-UTC", "true")
    }
    return req, err, url
}

func connect(requestType string,requestUrl string, jsonData io.Reader) []byte{
	checkAnnotation()
	var client = checkSSL()
	
    req, err, url := buildRequest(requestType,requestUrl,jsonData)
    
    rs, err := client.Do(req)

    if err != nil {
        fmt.Println(err)
        fmt.Println("\033[31mError: Could not connect to site '" + url + "'")
        os.Exit(1)
    }
    defer rs.Body.Close()
 
    if (requestType == "PATCH" || requestType == "DELETE") && rs.StatusCode == 200 {
        os.Exit(0)
    }

    bodyBytes, err := ioutil.ReadAll(rs.Body)
    
    if err != nil {
        panic(err)
    }

    return bodyBytes
}