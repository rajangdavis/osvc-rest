package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "osvc-rest",
	Short: "OSvC REST CLI",
	Long: `osvc-rest - a Command Line Interface application to work with the Oracle Service Cloud REST API

osvc-rest  Copyright (C) 2018 Rajan G. Davis
	
This program comes with ABSOLUTELY NO WARRANTY; for details type 'osvc-rest warranty'.
This is free software, and you are welcome to redistribute it
under certain conditions`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&userName, "username", "u", "", "Username to use for basic authentication")
	RootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password to use for basic authentication")
	RootCmd.PersistentFlags().StringVarP(&session, "session", "s", "", "Sets the session token for session authentication")
	RootCmd.PersistentFlags().StringVarP(&oauth, "oauth", "o", "", "Sets the OAuth token for OAuth authentication")
	RootCmd.PersistentFlags().StringVarP(&interfaceName, "interface", "i", "", "Oracle Service Cloud Interface to connect with")
	RootCmd.PersistentFlags().StringVarP(&vhost, "vhost", "", "", "Oracle Service Cloud virtual host to connect with")

	RootCmd.PersistentFlags().BoolVarP(&demoSite, "demosite", "", false, "Change the domain from 'custhelp' to 'rightnowdemo'")
	RootCmd.PersistentFlags().BoolVarP(&suppressRules, "suppress-rules", "", false, "Adds a header to suppress business rules")
	RootCmd.PersistentFlags().BoolVarP(&noSslVerify, "no-ssl-verify", "", false, "Turns off SSL verification")
	RootCmd.PersistentFlags().StringVarP(&version, "version", "v", "v1.3", "Changes the CCOM version")
	RootCmd.PersistentFlags().StringVarP(&accessToken, "access-token", "", "", "Adds access token header")
	RootCmd.PersistentFlags().StringVarP(&annotation, "annotate", "a", "", "Adds a custom header that adds an annotation (CCOM version must be set to \"v1.4\" or \"latest\"); limited to 40 characters")
	RootCmd.PersistentFlags().BoolVarP(&excludeNull, "exclude-null", "e", false, "Adds a custom header to excludes null from results")
	RootCmd.PersistentFlags().BoolVarP(&utcTime, "utc-time", "t", false, "Adds a custom header to return results using Coordinated Universal Time (UTC) format for time (Supported on November 2016+)")
	RootCmd.PersistentFlags().BoolVarP(&schema, "schema", "", false, "Sets 'Accept' header to 'application/schema+json'")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Prints request headers for debugging")
	RootCmd.PersistentFlags().IntVarP(&nextRequest, "next-request", "", 0, "Number of milliseconds before another HTTP request can be made with the associated access-token; this is an anti-DDoS measure")
}
