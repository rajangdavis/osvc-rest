package cmd

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/buger/jsonparser"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// For Uploading File Attachments
func openFile(fileAttachmentPath string) string {
	// Get the file name
	fileName := filepath.Base(fileAttachmentPath)

	// Open file on disk.
	f, fileError := os.Open(fileAttachmentPath)

	if fileError != nil {
		fmt.Println(fileName)
		fmt.Println("\033[31mError: there is something wrong with the file location you have input")
		os.Exit(0)
	}

	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)
	return fmt.Sprintf(`{"fileName" : "%s", "data" : "%s"}`, fileName, encoded)
}

// For Downloading File Attachments
func createFile(fileAttachmentPath string, body []byte) error {
	out, err := os.Create(fileAttachmentPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}

func downloadFileData(url string) string {
	newUrl := strings.Replace(url, "?download", "", 1)
	downloadData := connect("GET", newUrl, nil)
	fileName, err := jsonparser.GetString(downloadData, "fileName")
	if err != nil {
		fileName = "downloadedAttachment.tgz"
	}
	return fileName
}
