package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"runtime"
	"strings"
)

type TicketStruct struct {
	TicketId      string
	DirectoryPath string
	MetaDataPath  string
	FilesToCreate []string
}

func (ts *TicketStruct) setTicketId(args []string, envVar envVarStruct) error {
	if len(args) == 1 {
		ts.TicketId = args[0]
	} else if envVar.exists {
		ts.TicketId = envVar.value
	} else {
		return errors.New("no ticket id set")
	}

	return nil
}

func (ts *TicketStruct) getTicketDirectory(ticketsPath string) string {
	return ticketsPath + "/" + ts.TicketId
}

func urlFormatValidator(urlFormat string) error {
	parsedURL, err := url.Parse(urlFormat)
	if err != nil {
		return err
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("URL scheme must be http or https")
	}

	if strings.Count(urlFormat, "@") != 1 {
		return errors.New("URL must contain exactly one '@' symbol")
	}

	return nil
}

func unMarshallMetadata(metaData []byte) (map[string]string, error) {
	var meta map[string]string
	if err := json.Unmarshal(metaData, &meta); err != nil {
		return nil, err
	}
	return meta, nil

}

func addValuesToMetaJson(description string, urlFormat string, ticketId string) string {
	url := strings.Replace(urlFormat, "@", ticketId, 1)
	return `{
		"description": "` + description + `",
		"url": "` + url + `"
}`
}

func createMetaJson(ticketDirectoryPath string, description string, urlFormat string, ticketId string) error {
	if err := createFileWithContent(ticketDirectoryPath+"/meta.json", addValuesToMetaJson(description, urlFormat, ticketId)); err != nil {
		return err
	}

	return nil
}

func openTicketDirectory(ticketPath string, envVar envVarStruct) error {
	var editor []string
	if envVar.exists {
		editor = strings.Split(envVar.value, " ")
	} else {
		switch runtime.GOOS {
		case "darwin":
			editor = append(editor, "open")
		case "linux":
			editor = append(editor, "xdg-open")
		case "windows":
			editor = append(editor, "explorer")
		default:
			return fmt.Errorf("unsupported platform")
		}
	}

	if err := openFile(editor, ticketPath); err != nil {
		return err
	}

	return nil
}