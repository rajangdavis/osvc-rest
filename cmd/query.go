package cmd

import (
	"strings"
	"fmt"
	"github.com/spf13/cobra"
)

var query = &cobra.Command{
	Use: "query",
	Short: "Runs up to one or more ROQL queries",
	Long: "Runs up to one or more ROQL queries",
	RunE: func(cmd *cobra.Command, args []string) error {
		for i := 0; i < len(args); i++ {
			splitArgs := strings.Split(args[i],":")
			fmt.Println(fmt.Sprintf("key is %s", splitArgs[0]))
			fmt.Println(fmt.Sprintf("query is %s",splitArgs[1]))
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(query)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// byeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// byeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}