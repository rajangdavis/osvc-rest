package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"errors"
)

func checkRequiredFlags(flags *pflag.FlagSet) error {
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


// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "osvc-rest",
	Short: "OSvC REST CLI",
	Long: `osvc-rest - a Command Line Interface application to work with the Oracle Service Cloud REST API`,
	PreRunE: func(cmd *cobra.Command, args []string) error {		
		return checkRequiredFlags(cmd.Flags())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&userName,"username","u","", "Username to use for basic authentication")
    RootCmd.MarkPersistentFlagRequired("username")
    RootCmd.PersistentFlags().StringVarP(&password,"password","p","", "Password to use for basic authentication")
    RootCmd.MarkPersistentFlagRequired("password")
    RootCmd.PersistentFlags().StringVarP(&interfaceName,"interface","i","", "Oracle Service Cloud Interface to connect with")
    RootCmd.MarkPersistentFlagRequired("interface")
    
    RootCmd.PersistentFlags().BoolVarP(&demoSite,"demosite","",false, "Change the domain from 'custhelp' to 'rightnowdemo'")
    RootCmd.PersistentFlags().BoolVarP(&suppressRules,"suppress-rules","",false, "Adds a header to suppress business rules")
    RootCmd.PersistentFlags().BoolVarP(&noSslVerify,"no-ssl-verify","",false, "Turns off SSL verification")
    RootCmd.PersistentFlags().StringVarP(&version,"version","v","v1.3", "Changes the CCOM version")
    RootCmd.PersistentFlags().StringVarP(&annotation,"annotate","a","", "Adds a custom header that adds an annotation (CCOM version must be set to \"v1.4\" or \"latest\"); limited to 40 characters")
    RootCmd.PersistentFlags().BoolVarP(&excludeNull,"exclude-null","e",false, "Adds a custom header to excludes null from results")
    RootCmd.PersistentFlags().BoolVarP(&utcTime,"utcTime","t",false, "Adds a custom header to return results using Coordinated Universal Time (UTC) format for time (Supported on November 2016+)")
 
    RootCmd.PersistentFlags().BoolVarP(&debug,"debug","",false, "Prints request headers for debugging")
}