package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
)

var (
	userName, password, interfaceName, version, annotation, accessToken       string
	noSslVerify, suppressRules, demoSite, excludeNull, utcTime, debug, schema bool
	nextRequest                                                               int
)

func checkAnnotation() error {
	if annotation != "" && len(annotation) > 40 {
		fmt.Println("\033[31mError: Annotation cannot be greater than 40 characters.")
		os.Exit(0)
	} else if (version == "v1.4" || version == "latest") && annotation == "" {
		fmt.Println("\033[31mError: An Annotation must be set when using CCOM version 1.4 (e.g. -a \"40 character annotation\")")
		os.Exit(0)
	}
	return nil
}

func setDomain() string {
	domain := ""
	if demoSite == true {
		domain = "rightnowdemo"
	} else {
		domain = "custhelp"
	}
	return domain
}

func checkSSL() *http.Client {
	var client = &http.Client{}
	if noSslVerify == true {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}
	return client
}
