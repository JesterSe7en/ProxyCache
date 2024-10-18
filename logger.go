package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	Logger      *log.Logger
	fileToWrite *os.File // Move file handle to a global variable for proper closure handling
)

func initLogger() error {
	year, month, day := time.Now().Date()
	logFilename := fmt.Sprintf("%04d-%02d-%02d.txt", year, month, day)

	// Check if log already exists, if not, create it
	if _, err := os.Stat(logFilename); os.IsNotExist(err) {
		fileToWrite, err = os.Create(logFilename) // Store file handle in a global var
		if err != nil {
			return fmt.Errorf("cannot create log file: %w", err)
		}
	} else {
		// Append to existing file
		fileToWrite, err = os.OpenFile(logFilename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("cannot open log file: %w", err)
		}
	}

	// Initialize the logger with the file handle
	Logger = log.New(fileToWrite, "", log.Ldate|log.Ltime)
	return nil
}

func closeLogger() {
	// Ensure the file is closed when the program exits
	if fileToWrite != nil {
		fileToWrite.Close()
	}
}

// LogFatal logs fatal errors and exits the program
func LogFatal(message string, err error) {
	Logger.SetPrefix("FATAL: ")
	if err != nil {
		Logger.Fatalf("%s: %v", message, err)
	} else {
		Logger.Fatalln(message)
	}
}

// LogError logs non-fatal errors
func LogError(message string, err error) {
	Logger.SetPrefix("ERROR: ")
	if err != nil {
		Logger.Fatalf("%s: %v", message, err)
	} else {
		Logger.Println(message)
	}
}

// LogWarn logs warnings
func LogWarn(message string) {
	Logger.SetPrefix("WARN: ")
	Logger.Println(message)
}

// LogInfo logs informational messages
func LogInfo(message string) {
	Logger.SetPrefix("INFO: ")
	Logger.Println(message)
}
