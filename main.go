package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/aminnairi/gouache/lib/fs"
	"github.com/aminnairi/gouache/lib/logger"
	"github.com/aminnairi/gouache/lib/number"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

type Options struct {
	WithBody    bool
	WithHeaders bool
	WithStatus  bool
}

type HTTPRequest struct {
	filePath string
	path     string
	method   string
	body     string
	host     string
}

func (httpRequest HTTPRequest) incomplete() bool {
	return len(strings.TrimSpace(httpRequest.method)) == 0 ||
		len(strings.TrimSpace(httpRequest.host)) == 0 ||
		len(strings.TrimSpace(httpRequest.path)) == 0
}

func (httpRequest HTTPRequest) invalidHostPrefix() bool {
	prefixes := []string{"http://", "https://"}

	for _, prefix := range prefixes {
		if !strings.HasPrefix(httpRequest.host, prefix) {
			return false
		}
	}

	return true
}

func (httpRequest HTTPRequest) invalidFilePathSuffix() bool {
	return !strings.HasSuffix(strings.TrimSpace(httpRequest.filePath), ".http")
}

func (httpRequest HTTPRequest) invalidMethod() bool {
	allowedMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	return !slices.Contains(allowedMethods, httpRequest.method)
}

func (httpRequest HTTPRequest) invalidPath() bool {
	return !strings.HasPrefix(strings.TrimSpace(httpRequest.path), "/")
}

func main() {
	options := Options{}

	rootCommand := &cobra.Command{
		Use:   "gouache",
		Short: "Create and send HTTP request from files",
		Long:  "Create HTTP requests from files and run them right from your terminal",
	}

	requestCommand := &cobra.Command{
		Use:   "request file.http",
		Short: "Send HTTP requests",
		Long:  "Send HTTP requests to the provided file or folder containing files for each one of your HTTP requests",
		// TODO: send a request in interactive mode if no argument is passed
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, arguments []string) {
			allowedMethods := []string{"GET", "POST", "PATCH", "DELETE", "PUT"}
			filePaths := []string{}

			if len(arguments) == 0 {
				logger.Fatal("No file or folder provided.")
			}

			for _, argument := range arguments {
				for filePath := range fs.Files(argument) {
					if !strings.HasSuffix(filePath, ".http") {
						logger.Warning("File", filePath, "ignored because it does not have the suffix .http")
						continue
					}

					filePaths = append(filePaths, filePath)
				}
			}

			for _, filePath := range filePaths {
				logger.Info("Sending request from file", filePath)

				stat, statError := os.Stat(filePath)

				if statError != nil {
					logger.Fatal("File", filePath, "does not exist or is not readable.")
				}

				if stat.IsDir() {
					logger.Fatal("Provided request should not be a directory, but rather a path to a file")
				}

				file, openError := os.Open(filePath)

				if openError != nil {
					logger.Fatal("Unable to open file:", openError)
				}

				scanner := bufio.NewScanner(file)

				if !scanner.Scan() {
					logger.Fatal("Expected a request line, got nothing.")
				}

				line := scanner.Text()
				parts := strings.Split(line, " ")

				if len(parts) != 3 {
					logger.Fatal("Invaid line encountered for line:", parts)
				}

				method := parts[0]
				validMethod := slices.Contains(allowedMethods, method)

				if !validMethod {
					logger.Fatal("Invalid method:", method, "expected one of the following:", strings.Join(allowedMethods, ", "))
				}

				path := parts[1]
				version := parts[2]

				if version != "HTTP/2" {
					logger.Fatal("HTTP version must be HTTP/2")
				}

				if !scanner.Scan() {
					logger.Fatal("Request must contain at least one header")
				}

				hostHeader := scanner.Text()
				hostHeaderParts := strings.Split(hostHeader, ": ")

				if len(hostHeaderParts) != 2 {
					logger.Fatal("Header must be in the following format: HeaderName: HeaderValue")
				}

				hostHeaderName := strings.Trim(hostHeaderParts[0], " ")
				hostHeaderValue := strings.Trim(hostHeaderParts[1], " ")

				if hostHeaderName != "Host" {
					logger.Fatal("First header must be the Host header")
				}

				isHeaderPrefixedWithHTTP := strings.HasPrefix(hostHeaderValue, "http://")
				isHeaderPrefixedWithHTTPS := strings.HasPrefix(hostHeaderValue, "https://")
				isHeaderCorrectlyPrefixed := isHeaderPrefixedWithHTTP || isHeaderPrefixedWithHTTPS

				if !isHeaderCorrectlyPrefixed {
					logger.Fatal("Host header value should starts with http:// or https://")
				}

				request, requestError := http.NewRequest(method, fmt.Sprint(hostHeaderValue, path), nil)

				isBody := false
				requestBody := ""

				for scanner.Scan() {
					if isBody {
						requestBody += scanner.Text()
						continue
					}

					header := scanner.Text()

					if len(strings.Trim(header, " ")) == 0 {
						isBody = true
						continue
					}

					headerParts := strings.Split(header, ": ")

					if len(headerParts) != 2 {
						logger.Fatal("Header must be in the following format: HeaderName: HeaderValue")
					}

					headerName := strings.Trim(headerParts[0], " ")
					headerValue := strings.Trim(headerParts[1], " ")

					request.Header.Add(headerName, headerValue)
				}

				client := &http.Client{}

				if requestError != nil {
					logger.Fatal("Unable to run request:", requestError)
				}

				response, clientError := client.Do(request)

				if clientError != nil {
					logger.Fatal("Error while running the request:", clientError)
				}

				if options.WithStatus {
					if number.IntBetween(100, 199, response.StatusCode) {
						logger.HTTPInformational(response.Status)
					} else if number.IntBetween(200, 299, response.StatusCode) {
						logger.HTTPSuccess(response.Status)
					} else if number.IntBetween(300, 399, response.StatusCode) {
						logger.HTTPRedirection(response.Status)
					} else if number.IntBetween(400, 499, response.StatusCode) {
						logger.HTTPClientError(response.Status)
					} else if number.IntBetween(500, 599, response.StatusCode) {
						logger.HTTPServerError(response.Status)
					} else {
						fmt.Println("HTTP/2", response.Status)
					}
				}

				if options.WithHeaders {
					headersTable := table.New().BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#01579B"))).Headers("Header", "Value")

					for headerName, headerValues := range response.Header {
						headersTable.Row(headerName, strings.Join(headerValues, ", "))
					}

					fmt.Println(headersTable)
				}

				if options.WithBody {
					responseBytes, responseError := io.ReadAll(response.Body)

					if responseError != nil {
						logger.Fatal("Unable to fetch the response body:", responseError)
					}

					fmt.Println(string(responseBytes))
				}

				closeError := file.Close()

				if closeError != nil {
					logger.Fatal("Error while closing file:", closeError)
				}
			}
		},
	}

	httpRequest := HTTPRequest{}

	generateCommand := &cobra.Command{
		Use:   "generate index.http",
		Short: "Generate a request",
		Long:  "Generate a request in the HTTP format",
		Args:  cobra.RangeArgs(0, 1),

		Run: func(cmd *cobra.Command, args []string) {
			httpRequest.filePath = args[0]

			if httpRequest.incomplete() {
				confirmation := false

				form := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().Title("Name of the file").Validate(func(value string) error {
							if strings.HasSuffix(value, ".http") {
								return nil
							}

							return errors.New("file should end in .http")
						}).Value(&httpRequest.filePath),
						huh.NewSelect[string]().Title("HTTP method").Options(
							huh.NewOption("GET", "GET"),
							huh.NewOption("POST", "POST"),
							huh.NewOption("GET", "GET"),
							huh.NewOption("PATCH", "PATCH"),
							huh.NewOption("DELETE", "DELETE"),
						).Value(&httpRequest.method),
						huh.NewInput().Title("Path for the HTTP request").Validate(func(value string) error {
							if strings.HasPrefix(value, "/") {
								return nil
							}

							return errors.New("path must start with /")
						}).Value(&httpRequest.path),
						huh.NewInput().Title("Host name").Validate(func(value string) error {
							if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
								return nil
							}

							return errors.New("must start with http:// or https://")
						}).Value(&httpRequest.host),
						huh.NewText().Title("Body of the HTTP request").Value(&httpRequest.body),
						huh.NewConfirm().Title("Save the request?").Affirmative("Save").Negative("Cancel").Value(&confirmation),
					),
				)

				if runError := form.Run(); runError != nil {
					logger.Fatal("Error while running the form:", runError)
				}

				if !confirmation {
					logger.Info("Okay understood, I won't create anything.")
					return
				}
			}

			if httpRequest.invalidHostPrefix() {
				logger.Fatal("Host must start with either http:// or https://, got", httpRequest.host)
			}

			if httpRequest.invalidFilePathSuffix() {
				logger.Fatal("File path must end with .http, got", httpRequest.filePath)
			}

			if httpRequest.invalidMethod() {
				logger.Fatal("Method is invalid and must be one of the following: GET, POST, PATCH, PUT, DELETE, got", httpRequest.method)
			}

			if httpRequest.invalidPath() {
				logger.Fatal("Path must start with a /, got ", httpRequest.path)
			}

			directoryPath := filepath.Dir(httpRequest.filePath)

			if !fs.FolderExists(directoryPath) {
				directoryCreationConfirmation := false

				directoryCreationForm := huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().Title(fmt.Sprintln("Directory", directoryPath, "does not exist")).Affirmative("Create").Negative("Cancel").Value(&directoryCreationConfirmation),
					),
				)

				if directoryCreationFormRunError := directoryCreationForm.Run(); directoryCreationFormRunError != nil {
					logger.Fatal("Failed to run the form directory:", directoryCreationFormRunError)
				}

				if !directoryCreationConfirmation {
					log.Fatal("Not creating any directory.")
				}

				if mkdirError := os.MkdirAll(directoryPath, 0o755); mkdirError != nil {
					logger.Fatal("failed to create folder:", mkdirError)
				}
			}

			if fs.FileExist(httpRequest.filePath) {
				overwriteExistingFile := false

				confirmForm := huh.NewForm(
					huh.NewGroup(
						huh.NewConfirm().Title("File already exists").Affirmative("Overwrite").Negative("Cancel").Value(&overwriteExistingFile),
					),
				)

				if confirmFormRunError := confirmForm.Run(); confirmFormRunError != nil {
					log.Fatal("Failed to confirm overwriting of file", httpRequest.filePath)
				}

				if !overwriteExistingFile {
					logger.Info("Okay, I won't overwrite the file.")
					return
				}
			}

			logger.Info("Okay, I'll create the file for you!")
			data := fmt.Sprintf("%s %s HTTP/2\nHost: %s\n", httpRequest.method, httpRequest.path, httpRequest.host)

			httpRequest.body = strings.TrimSpace(httpRequest.body)

			if len(httpRequest.body) != 0 {
				data += fmt.Sprintln("")
				data += fmt.Sprintln(httpRequest.body)
			}

			if writeError := os.WriteFile(httpRequest.filePath, []byte(data), 0o644); writeError != nil {
				logger.Fatal("i can't write the file", httpRequest.filePath, "because:", writeError)
			}

			logger.Info("wrote new request to file", httpRequest.filePath)
		},
	}

	rootCommand.AddCommand(requestCommand)
	rootCommand.AddCommand(generateCommand)

	requestCommand.Flags().BoolVarP(&options.WithBody, "with-body", "b", false, "Display the raw body of the response")
	requestCommand.Flags().BoolVarP(&options.WithStatus, "with-status", "s", false, "Display the status line of the response")
	requestCommand.Flags().BoolVarP(&options.WithHeaders, "with-headers", "H", false, "Display the headers of the response")

	generateCommand.Flags().StringVarP(&httpRequest.method, "method", "m", "", "HTTP method, either GET, POST, PUT, PATCH or DELETE")
	generateCommand.Flags().StringVarP(&httpRequest.host, "host", "H", "", "Value for the Host header, must start with either http:// or https://")
	generateCommand.Flags().StringVarP(&httpRequest.body, "body", "b", "", "Body for the HTTP request")
	generateCommand.Flags().StringVarP(&httpRequest.path, "path", "p", "", "Path for the HTTP request")

	if commandError := rootCommand.Execute(); commandError != nil {
		logger.Fatal("Error when executing the command:", commandError)
	}
}
