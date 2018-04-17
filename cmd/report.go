package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var lookupName, id, filters string


func checkReportFlags(flags *pflag.FlagSet) error {

	flags.VisitAll(func(flag *pflag.Flag) {
		fmt.Println(flag.Name)
	})

	return nil
}


func runReport(cmd *cobra.Command, args []string) error {

	reportUrl := "analyticsReportResults"

	bodyBytes := connect("POST",reportUrl)

	// parses JSON and outputs to console
	normalize(bodyBytes)

	return nil
}


// report represents the report command
var report = &cobra.Command{
	Use:   "report",
	Short: "Runs an analytics report command",
	Long:  "Runs an analytics report command \nExample: \n$ osvc-rest",
	PreRunE:func(cmd *cobra.Command, args []string) error {		
		return checkReportFlags(cmd.Flags())
	},
	RunE: runReport,
}

func init() {
	RootCmd.AddCommand(report)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// report.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	report.Flags().StringVarP(&filters,"filters","f","", "Adds a custom header that adds an annotation")
	report.Flags().StringVarP(&lookupName,"name","n","", "Adds a custom header that adds an annotation")
	report.Flags().StringVarP(&id, "id", "","", "Adds a custom header that adds an annotation")
}
