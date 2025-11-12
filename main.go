package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	RequestFilePath string `short:"r" long:"request" description:"HTTP request file to run"`
	WithBody        bool   `short:"b" long:"with-body" description:"Display the raw body of the response"`
	WithHeaders     bool   `short:"H" long:"with-headers" description:"Display the headers of the response"`
	WithStatus      bool   `short:"s" long:"with-status" description:"Display the status line of the response"`
}

func main() {
	allowedMethods := []string{"GET", "POST", "PATCH", "DELETE", "PUT"}
	options := Options{}
	arguments, parseError := flags.Parse(&options)

	if parseError != nil {
		log.Fatal("Failed to parse arguments:", parseError)
	}

	if len(arguments) > 0 {
		log.Fatal("This program does not expect any arguments.")
	}

	stat, statError := os.Stat(options.RequestFilePath)

	if statError != nil {
		log.Fatal("Provided path is not a file")
	}

	if stat.IsDir() {
		log.Fatal("Provided request should not be a directory, but rather a path to a file")
	}

	file, openError := os.Open(options.RequestFilePath)

	if openError != nil {
		log.Fatal("Unable to open file:", openError)
	}

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		log.Fatal("Expected a request line, got nothing.")
	}

	line := scanner.Text()
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		log.Fatal("Invaid line encountered for line:", parts)
	}

	method := parts[0]
	validMethod := slices.Contains(allowedMethods, method)

	if !validMethod {
		log.Fatal("Invalid method:", method)
	}

	path := parts[1]
	version := parts[2]

	if version != "HTTP/2" {
		log.Fatal("HTTP version must be HTTP/2")
	}

	if !scanner.Scan() {
		log.Fatal("Request must contain at least one header")
	}

	hostHeader := scanner.Text()
	hostHeaderParts := strings.Split(hostHeader, ": ")

	if len(hostHeaderParts) != 2 {
		log.Fatal("Header must be in the following format: HeaderName: HeaderValue")
	}

	hostHeaderName := strings.Trim(hostHeaderParts[0], " ")
	hostHeaderValue := strings.Trim(hostHeaderParts[1], " ")

	if hostHeaderName != "Host" {
		log.Fatal("First header must be the Host header")
	}

	isHeaderPrefixedWithHTTP := strings.HasPrefix(hostHeaderValue, "http://")
	isHeaderPrefixedWithHTTPS := strings.HasPrefix(hostHeaderValue, "https://")
	isHeaderCorrectlyPrefixed := isHeaderPrefixedWithHTTP || isHeaderPrefixedWithHTTPS

	if !isHeaderCorrectlyPrefixed {
		log.Fatal("Host header value should starts with http:// or https://")
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
			log.Fatal("Header must be in the following format: HeaderName: HeaderValue")
		}

		headerName := strings.Trim(headerParts[0], " ")
		headerValue := strings.Trim(headerParts[1], " ")

		request.Header.Add(headerName, headerValue)
	}

	client := &http.Client{}

	if requestError != nil {
		log.Fatal("Unable to run request:", requestError)
	}

	response, clientError := client.Do(request)

	if clientError != nil {
		log.Fatal("Error while running the request:", clientError)
	}

	if options.WithStatus {
		fmt.Println("HTTP/2", response.Status)
	}

	if options.WithHeaders {
		for headerName, headerValues := range response.Header {
			headerLine := fmt.Sprintf("%s: ", headerName)

			for _, headerValue := range headerValues {
				headerLine = fmt.Sprint(headerLine, headerValue)
			}

			fmt.Println(headerLine)
		}
	}

	if options.WithBody {
		responseBytes, responseError := io.ReadAll(response.Body)

		if responseError != nil {
			log.Fatal("Unable to fetch the response body:", responseError)
		}

		fmt.Println(string(responseBytes))
	}

	closeError := file.Close()

	if closeError != nil {
		log.Fatal("Error while closing file:", closeError)
	}
}
