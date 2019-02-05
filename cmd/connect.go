package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
)

func buildRequest(method string, requestUrl string, jsonData io.Reader) (*http.Request, error, string) {
	var baseUrl = "https://" + setDomain() + "/services/rest/connect/" + version + "/"
	var finalUrl = baseUrl + strings.Replace(requestUrl, " ", "%20", -1)

	var req *http.Request
	var err error

	if method == "POST" || method == "PATCH" {
		req, err = http.NewRequest("POST", finalUrl, jsonData)
		req.Header.Set("Content-Type", "application/json")
		if method == "PATCH" {
			req.Header.Set("X-HTTP-Method-Override", "PATCH")
		}
	} else {
		req, err = http.NewRequest(method, finalUrl, nil)
	}

	if session != "" {
		req.Header.Add("Authorization", "Session "+session)
	} else if oauth != "" {
		req.Header.Add("Authorization", "Bearer "+oauth)
	} else {
		req.Header.Add("Authorization", "Basic "+basicAuth(userName, password))
	}

	if (version == "v1.4" || version == "latest") && annotation != "" && len(annotation) <= 40 {
		req.Header.Add("OSvC-CREST-Application-Context", annotation)
	}

	if excludeNull == true {
		req.Header.Add("prefer", "exclude-null-properties")
	}

	if utcTime == true {
		req.Header.Add("OSvC-CREST-Time-UTC", "yes")
	}

	if schema == true {
		req.Header.Add("Accept", "application/schema+json")
	}

	if accessToken != "" {
		req.Header.Add("osvc-crest-api-access-token", accessToken)
	}

	if nextRequest > 0 {
		req.Header.Add("osvc-crest-next-request-after", strconv.Itoa(nextRequest))
	}

	if debug == true {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(requestDump))
	}

	return req, err, baseUrl
}

func connect(requestType string, requestUrl string, jsonData io.Reader) []byte {
	checkAnnotation()
	var client = checkSSL()

	req, err, url := buildRequest(requestType, requestUrl, jsonData)

	rs, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		fmt.Println("\033[31mError: Could not connect to site '" + url + "' \033[0m ")
		os.Exit(1)
	}
	defer rs.Body.Close()

	if (requestType == "PATCH" || requestType == "DELETE") && rs.StatusCode == 200 {
		os.Exit(0)
	}

	var bodyBytes []byte

	if requestType == "OPTIONS" {
		responseDump, err := httputil.DumpResponse(rs, true)
		if err != nil {
			fmt.Println(err)
		}
		bodyBytes = responseDump
	} else {
		bodyBytes, err = ioutil.ReadAll(rs.Body)
	}

	if err != nil {
		panic(err)
	}

	return bodyBytes
}
