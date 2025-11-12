// Package logger gathers all of the functions useful for logging errors and messages to the console
package logger

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func Fatal(messages ...any) {
	errorStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#B00020")).Foreground(lipgloss.Color("#FFFFFF")).PaddingLeft(1).PaddingRight(1)
	allMessages := []any{errorStyle.Render("ERROR")}
	allMessages = append(allMessages, messages...)

	fmt.Println(allMessages...)
	os.Exit(1)
}

func Warning(messages ...any) {
	errorStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#EF6C00")).Foreground(lipgloss.Color("#FFFFFF")).PaddingLeft(1).PaddingRight(1)
	allMessages := []any{errorStyle.Render("WARNING")}
	allMessages = append(allMessages, messages...)

	fmt.Println(allMessages...)
}

func Info(messages ...any) {
	errorStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#01579B")).Foreground(lipgloss.Color("#FFFFFF")).PaddingLeft(1).PaddingRight(1)
	allMessages := []any{errorStyle.Render("INFO")}
	allMessages = append(allMessages, messages...)

	fmt.Println(allMessages...)
}

func HTTPInformational(statusText string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#FFFFFF")).Bold(true).PaddingLeft(1).PaddingRight(1)
	fmt.Println("HTTP/2", style.Render(statusText))
}

func HTTPSuccess(statusText string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#1B5E20")).Bold(true).PaddingLeft(1).PaddingRight(1)
	fmt.Println("HTTP/2", style.Render(statusText))
}

func HTTPRedirection(statusText string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#01579B")).Bold(true).PaddingLeft(1).PaddingRight(1)
	fmt.Println("HTTP/2", style.Render(statusText))
}

func HTTPClientError(statusText string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#E651000")).Bold(true).PaddingLeft(1).PaddingRight(1)
	fmt.Println("HTTP/2", style.Render(statusText))
}

func HTTPServerError(statusText string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#B71C1C")).Bold(true).PaddingLeft(1).PaddingRight(1)
	fmt.Println("HTTP/2", style.Render(statusText))
}
