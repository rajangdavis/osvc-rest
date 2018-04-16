package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// report represents the report command
var report = &cobra.Command{
	Use:   "report",
	Short: "Runs an analytics report command",
	Long:  "Runs an analytics report command \nExample: \n$ osvc-rest",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bye called")
	},
}

func init() {
	RootCmd.AddCommand(report)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// report.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// report.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
