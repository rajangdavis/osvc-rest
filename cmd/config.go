package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
)

var (
	userName, password, interfaceName, session, oauth, version, annotation, accessToken string
	noSslVerify, suppressRules, demoSite, excludeNull, utcTime, debug, schema           bool
	nextRequest                                                                         int
)

func interfaceAndPassword() error {
	checkInterface()
	checkAuthentication()
	return nil
}

func checkInterface() error {
	if interfaceName == "" {
		fmt.Println("\033[31mError: Must set an interface to connect with.")
		os.Exit(0)
	}
	return nil
}

func checkAuthentication() error {
	if userName == "" && password != "" {
		fmt.Println("\033[31mError: Password is set but user name is not")
		os.Exit(0)
	} else if userName != "" && password == "" {
		fmt.Println("\033[31mError: User name is set but password is not.")
		os.Exit(0)
	} else if userName == "" && password == "" && session == "" && oauth == "" {
		fmt.Println("\033[31mError: Must use some form of authentication.")
		os.Exit(0)
	}
	return nil
}

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
