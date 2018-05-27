package cmd

import (
	"fmt"
	"os"
	"bytes"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var lookupName, filters string
var id, reportLimit int

func checkReportFlags(flags *pflag.FlagSet) error {

	if lookupName == "" && id == 0 {
		fmt.Println("\033[31mError: Must use either the --name or --id flag with for working with the AnalyticsReportResults object")
		os.Exit(0)
	}

	return nil
}

func runReport(cmd *cobra.Command, args []string) error {

	var identifier  []byte
	var str string

	if lookupName =="" { 
		str = fmt.Sprintf(`{"id":%d}`, id)
	}else{ 
		str = fmt.Sprintf(`{"lookupName":%q}`, lookupName)
	}
	

	if reportLimit > 0 {
		str = str[:len(str) - 1] + fmt.Sprintf(`, "limit" : %d}`,reportLimit)
	}

	if filters != ""{
		str = str[:len(str) - 1] + fmt.Sprintf(`, "filters" : %s}`, filters)
	}

	identifier = []byte(str)

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
	Long: "\033[93mRuns an analytics report and returns parsed results\033[0m \033[0;32m\n\nReport (without filters) Example: \033[0m \n$ osvc-rest report --id 176 -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE \033[0;32m\n\nReport (with filters and limiting) Example: \033[0m \n$ osvc-rest report --id 176 --limit 10 --filters '[{\"name\":\"search_ex\",\"values\":\"returns\"}]' -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE",
	PreRunE:func(cmd *cobra.Command, args []string) error {		
		return checkReportFlags(cmd.Flags())
	},
	RunE: runReport,
}

func init() {
	report.Flags().StringVarP(&filters,"filters","f","", "Adds filters for reporting")
	report.Flags().StringVarP(&lookupName,"name","n","", "Sets the lookupName of the AnalyticsReport that we wish to run")
	report.Flags().IntVarP(&reportLimit,"limit","l",0, "Adds limit for reporting")
	report.Flags().IntVarP(&id, "id", "",0, "Sets the id of the AnalyticsReport that we wish to run")
	RootCmd.AddCommand(report)
}