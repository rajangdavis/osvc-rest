package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "osvc-rest",
	Short: "OSvC REST CLI",
	Long: `osvc-rest - a Command Line Interface application to work with the Oracle Service Cloud REST API`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
