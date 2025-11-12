package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/aminnairi/gouache/lib/fs"
	"github.com/aminnairi/gouache/lib/logger"
	"github.com/aminnairi/gouache/lib/number"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jessevdk/go-flags"
)

// TODO: replace go-flags with https://github.com/spf13/cobra
// TODO: create a command for running requests
// TODO: create a command for generating requests interactively using https://github.com/charmbracelet/huh

type Options struct {
	WithBody    bool `short:"b" long:"with-body" description:"Display the raw body of the response"`
	WithHeaders bool `short:"H" long:"with-headers" description:"Display the headers of the response"`
	WithStatus  bool `short:"s" long:"with-status" description:"Display the status line of the response"`
}

func main() {
	allowedMethods := []string{"GET", "POST", "PATCH", "DELETE", "PUT"}
	options := Options{}
	arguments, parseError := flags.Parse(&options)
	filePaths := []string{}

	if parseError != nil {
		logger.Fatal("Failed to parse options:", parseError)
	}

	if len(arguments) == 0 {
		logger.Fatal("No file or folder provided.")
	}

	for _, argument := range arguments {
		for filePath := range fs.Files(argument) {
			if !strings.HasSuffix(filePath, ".http") {
				logger.Info("File", filePath, "ignored because it does not have the suffix .http")
				continue
			}

			filePaths = append(filePaths, filePath)
		}
	}

	for _, filePath := range filePaths {
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
}
