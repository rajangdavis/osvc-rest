package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tidwall/gjson"
	"io"
	"os"
	"strings"
)

func httpCheck(args []string) []string {
	interfaceAndPassword()
	resourceUrls := []string{}
	if len(args) == 0 {
		fmt.Println("\033[31mError: Must set at least one resource url \033[0m ")
		os.Exit(0)
	} else {
		for i := 0; i < len(args); i++ {
			resourceUrls = append(resourceUrls, args[i])

		}
	}
	return resourceUrls
}

func makeRequest(verb string, url string, optionalJson io.Reader, ch chan<- []byte) {
	byteData := connect(verb, url, optionalJson)
	m, ok := gjson.Parse(string(byteData)).Value().(map[string]interface{})
	jsonData := []byte{}

	if !ok && strings.Index(url, "?download") != -1 {
		fileName := downloadFileData(url)
		createFile(fileName, byteData)
	} else if !ok && verb == "OPTIONS" {
		jsonData = byteData
	} else if !ok {
		fmt.Println("Error")
	} else {
		formattedJson, _ := json.MarshalIndent(m, "", "  ")
		jsonData = formattedJson
	}

	ch <- jsonData

}

func runHttp(cmd *cobra.Command, args []string) error {

	resourceUrls := httpCheck(args)
	resourceUrlsCount := len(resourceUrls)
	httpVerb := strings.ToUpper(cmd.Use)

	ch := make(chan []byte)

	if len(fileAttachmentsLocation) > 0 {

		if data == "" {
			data = `{"fileAttachments" : [ `
		} else {
			data = data[:len(data)-1] + `, "fileAttachments" : [ `

		}

		for i := 0; i < len(fileAttachmentsLocation); i++ {
			fileAttachmentData := openFile(fileAttachmentsLocation[i])

			if i < len(fileAttachmentsLocation)-1 {
				data = data + fileAttachmentData + ", "
			} else {
				data = data + fileAttachmentData
			}
		}

		data = data + `]}`

	}

	jsonData := bytes.NewReader([]byte(data))

	for i := 0; i < resourceUrlsCount; i++ {
		go makeRequest(httpVerb, resourceUrls[i], jsonData, ch)

		fmt.Fprintf(os.Stdout, "%s", <-ch)
		if httpVerb != "GET" {
			return nil
		}
	}

	return nil
}

var get = &cobra.Command{
	Use:   "get",
	Short: "Performs one or more GET requests",
	Long:  "\033[93mPerforms one or more GET requests and returns parsed results\033[0m \033[0;32m\n\nSingle Query Example: \033[0m \n$ osvc-rest query \"DESCRIBE\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE \033[0;32m\n\nMultiple Queries Example:\033[0m \n$ osvc-rest query \"SELECT * FROM INCIDENTS LIMIT 100\" \"SELECT * FROM SERVICEPRODUCTS LIMIT 100\" \"SELECT * FROM SERVICECATEGORIES LIMIT 100\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE",
	RunE:  runHttp,
}

var delete = &cobra.Command{
	Use:   "delete",
	Short: "Performs a DELETE request",
	Long:  "\033[93mPerforms a DELETE request; if successful, nothing is returned \033[0m \033[0;32m\n\nExample: \033[0m \n$ osvc-rest delete \"opportunities/1\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE\n\n",
	RunE:  runHttp,
}

var options = &cobra.Command{
	Use:   "options",
	Short: "Performs a OPTIONS request",
	Long:  "\033[93mPerforms a OPTIONS request; if successful, HEADERS are returned \033[0m \033[0;32m\n\nExample: \033[0m \n$ osvc-rest options \"opportunities\" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE\n\n",
	RunE:  runHttp,
}

var data string
var fileAttachmentsLocation []string

func checkPostPatchFlags(flags *pflag.FlagSet) error {

	// Don't force people to have to set the data flag if they only want to attach a file
	// Otherwise if there is no file attachmeng AND no data
	// Raise a complaint

	if data == "" && len(fileAttachmentsLocation) == 0 {
		fmt.Println("\033[31mError: Must send JSON Data for POST and PATCH requests; use the --data flag \033[0m")
		os.Exit(0)
	}

	return nil
}

var post = &cobra.Command{
	Use:   "post",
	Short: "Performs a POST request",
	Long:  "\033[93mPerforms a POST request and returns parsed results\033[0m \033[0;32m\n\nExample: \033[0m \n$ osvc-rest post \"opportunities\" --data '{\"name\":\"PCS- 100 laptops\"}' -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE\n\n",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkPostPatchFlags(cmd.Flags())
	},
	RunE: runHttp,
}

var patch = &cobra.Command{
	Use:   "patch",
	Short: "Performs a PATCH request",
	Long:  "\033[93mPerforms a PATCH request; if successful, nothing is returned \033[0m \033[0;32m\n\nExample: \033[0m \n$ osvc-rest patch \"opportunities/1\" --data '{\"name\":\"PCS- 100 laptops UPDATED\"}' -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE\n\n",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkPostPatchFlags(cmd.Flags())
	},
	RunE: runHttp,
}

func init() {
	RootCmd.AddCommand(get)
	post.Flags().StringVarP(&data, "data", "j", "", "Sets the JSON data to be sent for the POST request")
	patch.Flags().StringVarP(&data, "data", "j", "", "Sets the JSON data to be sent for the PATCH request")

	post.Flags().StringArrayVarP(&fileAttachmentsLocation, "attach-file", "f", []string{}, "Sets the File location of the file attachment to be included with a POST request")
	patch.Flags().StringArrayVarP(&fileAttachmentsLocation, "attach-file", "f", []string{}, "Sets the File location of the file attachment to be included with a PATCH request")

	RootCmd.AddCommand(post)
	RootCmd.AddCommand(patch)
	RootCmd.AddCommand(delete)
	RootCmd.AddCommand(options)
}
