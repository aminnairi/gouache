package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

func main() {
	requestPath := flag.String("request", "", "gouache --request request.http")
	allowedMethods := []string{"GET", "POST", "PATCH", "DELETE", "PUT"}

	flag.Parse()

	stat, statError := os.Stat(*requestPath)

	if statError != nil {
		log.Fatal("Provided path is not a file")
	}

	if stat.IsDir() {
		log.Fatal("Provided request should not be a directory, but rather a path to a file")
	}

	file, openError := os.Open("requests/index.http")

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

	header := scanner.Text()
	headerParts := strings.Split(header, ": ")

	if len(headerParts) != 2 {
		log.Fatal("Header must be in the following format: HeaderName: HeaderValue")
	}

	headerName := strings.Trim(headerParts[0], " ")
	headerValue := strings.Trim(headerParts[1], " ")

	if headerName != "Host" {
		log.Fatal("First header must be the Host header")
	}

	request, requestError := http.NewRequest(method, fmt.Sprint(headerValue, path), nil)

	isBody := false
	body := ""

	for scanner.Scan() {
		if isBody {
			body += scanner.Text()
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

	fmt.Println("HTTP/2", response.Status)

	for headerName, headerValue := range response.Header {
		fmt.Printf("%s: %s\n", headerName, headerValue)
	}

	closeError := file.Close()

	if closeError != nil {
		log.Fatal("Error while closing file:", closeError)
	}
}
