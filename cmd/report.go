package cmd

import (
	"bytes"
	"fmt"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
)

var lookupName, filters, csvName string
var id, reportLimit, reportOffset int

func checkReportFlags(flags *pflag.FlagSet) error {
	interfaceAndPassword()
	if lookupName == "" && id == 0 {
		fmt.Println("\033[31mError: Must use either the --name or --id flag with for working with the AnalyticsReportResults object\033[0m ")
		os.Exit(0)
	}

	return nil
}

func runReport(cmd *cobra.Command, args []string) error {

	var identifier []byte
	var str string

	if lookupName == "" {
		str = fmt.Sprintf(`{"id":%d}`, id)
	} else {
		str = fmt.Sprintf(`{"lookupName":%q}`, lookupName)
	}

	if reportLimit > 0 {
		str = str[:len(str)-1] + fmt.Sprintf(`, "limit" : %d}`, reportLimit)
	}

	if reportOffset > 0 {
		str = str[:len(str)-1] + fmt.Sprintf(`, "offset" : %d}`, reportOffset)
	}

	if filters != "" {
		str = str[:len(str)-1] + fmt.Sprintf(`, "filters" : %s}`, filters)
	}

	identifier = []byte(str)

	jsonData := bytes.NewBuffer(identifier)

	var results []map[string]interface{}
	bodyBytes := connect("POST","analyticsReportResults", jsonData)

	if csvName != ""{
		csvReport(bodyBytes, csvName)
	}else{
		finalResults := normalizeReport(bodyBytes, identifier, &results)
		jsonDataFinal, _ := json.MarshalIndent(finalResults, "", "  ")
		fmt.Fprintf(os.Stdout, "%s", jsonDataFinal)
	}
	return nil
}

// report represents the report command
var report = &cobra.Command{
	Use:   "report",
	Short: "Runs an analytics report command",
	Long:  "\033[93mRuns an analytics report and returns parsed results\033[0m \033[0;32m\n\nReport (without filters) Example: \033[0m \n$ osvc-rest report --id 176 -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE \033[0;32m\n\nReport (with filters, offsetting and limiting results) Example: \033[0m \n$ osvc-rest report --id 176 --limit 10 --offset 10 --filters '[{\"name\":\"search_ex\",\"values\":\"returns\"}]' -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkReportFlags(cmd.Flags())
	},
	RunE: runReport,
}

func init() {
	report.Flags().StringVarP(&filters, "filters", "f", "", "Adds filters for reporting")
	report.Flags().StringVarP(&lookupName, "name", "n", "", "Sets the lookupName of the AnalyticsReport that we wish to run")
	report.Flags().StringVarP(&csvName, "csv", "", "", "Exports to CSV to the file name provided")
	report.Flags().IntVarP(&reportLimit, "limit", "l", 0, "Adds limit for reporting")
	report.Flags().IntVarP(&reportOffset, "offset", "", 0, "Adds and offset for reporting")
	report.Flags().IntVarP(&id, "id", "", 0, "Sets the id of the AnalyticsReport that we wish to run")
	RootCmd.AddCommand(report)
}
