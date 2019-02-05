package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
)

var (
	userName, password, interfaceName, session, oauth, version, annotation, accessToken, vhost  string
	noSslVerify, suppressRules, demoSite, excludeNull, utcTime, debug, schema, keepAlive bool
	nextRequest                                                                          int
)

func interfaceAndPassword() error {
	checkInterface()
	checkAuthentication()
	return nil
}

func checkInterface() error {
	if interfaceName == ""  && vhost == ""{
		fmt.Println("\033[31mError: Must set an interface or vhost to connect with. \033[0m ")
		os.Exit(0)
	}
	return nil
}

func checkAuthentication() error {
	if userName == "" && password != "" {
		fmt.Println("\033[31mError: Password is set but user name is not \033[0m ")
		os.Exit(0)
	} else if userName != "" && password == "" {
		fmt.Println("\033[31mError: User name is set but password is not. \033[0m ")
		os.Exit(0)
	} else if userName == "" && password == "" && session == "" && oauth == "" {
		fmt.Println("\033[31mError: Must use some form of authentication. \033[0m ")
		os.Exit(0)
	}
	return nil
}

func checkAnnotation() error {
	if annotation != "" && len(annotation) > 40 {
		fmt.Println("\033[31mError: Annotation cannot be greater than 40 characters. \033[0m ")
		os.Exit(0)
	} else if (version == "v1.4" || version == "latest") && annotation == "" {
		fmt.Println("\033[31mError: An Annotation must be set when using CCOM version 1.4 (e.g. -a \"40 character annotation\") \033[0m ")
		os.Exit(0)
	}
	return nil
}

func setDomain() string {
	domain := ""
	if demoSite == true {
		domain = interfaceName + ".rightnowdemo.com"
	}else if vhost != ""{
		domain = vhost
	}else {
		domain = interfaceName + ".custhelp.com"
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
