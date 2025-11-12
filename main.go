package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	RequestFilePath string `short:"r" long:"request" description:"HTTP request file to run"`
	WithBody        bool   `short:"b" long:"with-body" description:"Display the raw body of the response"`
	WithHeaders     bool   `short:"H" long:"with-headers" description:"Display the headers of the response"`
	WithStatus      bool   `short:"s" long:"with-status" description:"Display the status line of the response"`
}

func errorAndExit(messages ...any) {
	errorStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#B00020")).Foreground(lipgloss.Color("#FFFFFF")).PaddingLeft(1).PaddingRight(1)
	allMessages := []any{errorStyle.Render("ERROR")}
	allMessages = append(allMessages, messages...)

	fmt.Println(allMessages...)
	os.Exit(1)
}

func main() {
	allowedMethods := []string{"GET", "POST", "PATCH", "DELETE", "PUT"}
	options := Options{}
	arguments, parseError := flags.Parse(&options)

	if parseError != nil {
		errorAndExit("Failed to parse options:", parseError)
	}

	if len(arguments) > 0 {
		errorAndExit("This program does not expect any arguments.")
	}

	stat, statError := os.Stat(options.RequestFilePath)

	if statError != nil {
		errorAndExit("File", options.RequestFilePath, "does not exist or is not readable.")
	}

	if stat.IsDir() {
		errorAndExit("Provided request should not be a directory, but rather a path to a file")
	}

	file, openError := os.Open(options.RequestFilePath)

	if openError != nil {
		errorAndExit("Unable to open file:", openError)
	}

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		errorAndExit("Expected a request line, got nothing.")
	}

	line := scanner.Text()
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		errorAndExit("Invaid line encountered for line:", parts)
	}

	method := parts[0]
	validMethod := slices.Contains(allowedMethods, method)

	if !validMethod {
		errorAndExit("Invalid method:", method)
	}

	path := parts[1]
	version := parts[2]

	if version != "HTTP/2" {
		errorAndExit("HTTP version must be HTTP/2")
	}

	if !scanner.Scan() {
		errorAndExit("Request must contain at least one header")
	}

	hostHeader := scanner.Text()
	hostHeaderParts := strings.Split(hostHeader, ": ")

	if len(hostHeaderParts) != 2 {
		errorAndExit("Header must be in the following format: HeaderName: HeaderValue")
	}

	hostHeaderName := strings.Trim(hostHeaderParts[0], " ")
	hostHeaderValue := strings.Trim(hostHeaderParts[1], " ")

	if hostHeaderName != "Host" {
		errorAndExit("First header must be the Host header")
	}

	isHeaderPrefixedWithHTTP := strings.HasPrefix(hostHeaderValue, "http://")
	isHeaderPrefixedWithHTTPS := strings.HasPrefix(hostHeaderValue, "https://")
	isHeaderCorrectlyPrefixed := isHeaderPrefixedWithHTTP || isHeaderPrefixedWithHTTPS

	if !isHeaderCorrectlyPrefixed {
		errorAndExit("Host header value should starts with http:// or https://")
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
			errorAndExit("Header must be in the following format: HeaderName: HeaderValue")
		}

		headerName := strings.Trim(headerParts[0], " ")
		headerValue := strings.Trim(headerParts[1], " ")

		request.Header.Add(headerName, headerValue)
	}

	client := &http.Client{}

	if requestError != nil {
		errorAndExit("Unable to run request:", requestError)
	}

	response, clientError := client.Do(request)

	if clientError != nil {
		errorAndExit("Error while running the request:", clientError)
	}

	if options.WithStatus {
		if response.StatusCode >= 100 && response.StatusCode <= 199 {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFFFF")).Bold(true).PaddingLeft(1).PaddingRight(1)
			fmt.Println("HTTP/2", style.Render(response.Status))
		} else if response.StatusCode >= 200 && response.StatusCode <= 299 {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#1B5E20")).Bold(true).PaddingLeft(1).PaddingRight(1)
			fmt.Println("HTTP/2", style.Render(response.Status))
		} else if response.StatusCode >= 300 && response.StatusCode <= 399 {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#01579B")).Bold(true).PaddingLeft(1).PaddingRight(1)
			fmt.Println("HTTP/2", style.Render(response.Status))
		} else if response.StatusCode >= 400 && response.StatusCode <= 499 {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#E651000")).Bold(true).PaddingLeft(1).PaddingRight(1)
			fmt.Println("HTTP/2", style.Render(response.Status))
		} else if response.StatusCode >= 500 && response.StatusCode <= 599 {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#B71C1C")).Bold(true).PaddingLeft(1).PaddingRight(1)
			fmt.Println("HTTP/2", style.Render(response.Status))
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
			errorAndExit("Unable to fetch the response body:", responseError)
		}

		fmt.Println(string(responseBytes))
	}

	closeError := file.Close()

	if closeError != nil {
		errorAndExit("Error while closing file:", closeError)
	}
}
