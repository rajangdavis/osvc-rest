package cmd

import (
	"fmt"
	"os"
	"bytes"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var lookupName, filters string
var id int

func checkReportFlags(flags *pflag.FlagSet) error {

	if lookupName == "" && id == 0 {
		fmt.Println("\033[31mError: Must use either the --name or --id flag with for working with the AnalyticsReportResults object")
		os.Exit(0)
	}

	return nil
}

func runReport(cmd *cobra.Command, args []string) error {

	var identifier  []byte

	if lookupName =="" { 
		str := fmt.Sprintf(`{"id":%d}`, id)
		identifier = []byte(str)
	}else{ 
		str := fmt.Sprintf(`{"lookupName":%q}`, lookupName)
		identifier = []byte(str)
	}
	jsonData := bytes.NewBuffer(identifier)
	reportUrl := "analyticsReportResults"
	bodyBytes := connect("POST",reportUrl,jsonData)
	normalizeReport(bodyBytes)
	return nil
}

// report represents the report command
var report = &cobra.Command{
	Use:   "report",
	Short: "Runs an analytics report command",
	Long: "\033[93mRuns an analytics report and returns parsed results\033[0m \033[0;32m\n\nReport (without filters) Example: \033[0m \n$ osvc-rest report --id 186 -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE",
	PreRunE:func(cmd *cobra.Command, args []string) error {		
		return checkReportFlags(cmd.Flags())
	},
	RunE: runReport,
}

func init() {
	RootCmd.AddCommand(report)
	report.Flags().StringVarP(&filters,"filters","f","", "Adds filters for reporting")
	report.Flags().StringVarP(&lookupName,"name","n","", "Sets the lookupName of the AnalyticsReport that we wish to run")
	report.Flags().IntVarP(&id, "id", "",0, "Sets the id of the AnalyticsReport that we wish to run")
	// report.Flags().IntVarP(&id, "id", "",0, "Adds a custom header that adds an annotation") add one for offset
	// report.Flags().IntVarP(&id, "id", "",0, "Adds a custom header that adds an annotation")
}