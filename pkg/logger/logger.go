package logger

import (
	"github.com/fatih/color"
	"time"
)

func Log(message string) {
	if message == "" {
		return
	}
	message = formatMessage(message, "Log", "    ")
	color.White(message)
}

func Info(message string) {
	if message == "" {
		return
	}
	message = formatMessage(message, "Info", "   ")
	color.White(message)
}

func Warn(message string) {
	if message == "" {
		return
	}
	message = formatMessage(message, "Warn", "   ")
	color.Yellow(message)
}

func Error(message string) {
	if message == "" {
		return
	}
	message = formatMessage(message, "Error", "  ")
	color.Red(message)
}

func Command(message string) {
	if message == "" {
		return
	}
	message = formatMessage(message, "Command", "")
	color.Cyan(message)
}

func formatMessage(message string, logLevel string, pad string) string {
	return getTimeStamp() + " [" + logLevel + "] " + pad + message
}

func getTimeStamp() string {
	return time.Now().Local().Format("2006/01/02 15:04:05")
}

