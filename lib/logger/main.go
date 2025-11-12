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

func Info(messages ...any) {
	errorStyle := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("#01579B")).Foreground(lipgloss.Color("#FFFFFF")).PaddingLeft(1).PaddingRight(1)
	allMessages := []any{errorStyle.Render("INFO")}
	allMessages = append(allMessages, messages...)

	fmt.Println(allMessages...)
}
