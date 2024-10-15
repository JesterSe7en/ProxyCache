package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var Logger *log.Logger

func initLogger() error {
	// Check if log already exist, if not create it

	todayDate, err := time.Parse("DateOnly", time.Now().String())
	if err != nil {
		return fmt.Errorf("cannot determine today's date")
	}

	var fileToWrite *os.File
	logFilename := fmt.Sprintf("%s.txt", todayDate.String())

	if _, err := os.Stat(logFilename); os.IsNotExist(err) {
		fileToWrite, err := os.Create(logFilename)
		if err != nil {
			defer fileToWrite.Close()
			return fmt.Errorf("cannot create log file: %w", err)
		}
	} else {
		// append to existing file
		fileToWrite, err := os.OpenFile(logFilename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			defer fileToWrite.Close()
			return fmt.Errorf("cannot open log file: %w", err)
		}
	}

	Logger = log.New(fileToWrite, "", log.Ldate|log.Ltime)
	return nil
}

// Will log an error and exit via os.Exit(1)
// Caution! Does not execute defer functions
func LogFatal(message string, err error) {
	Logger.SetPrefix("FATAL: ")
	if err != nil {
		Logger.Fatalln(fmt.Sprintf("%s: %v", message, err))
	} else {
		Logger.Fatalln(message)
	}
}

func LogError(message string, err error) {
	Logger.SetPrefix("ERROR: ")
	if err != nil {
		Logger.Println(fmt.Sprintf("%s: %v", message, err))
	} else {
		Logger.Println(message)
	}
}

func LogWarn(message string) {
	Logger.SetPrefix("WARN: ")
	Logger.Println(message)
}

func LogInfo(message string) {
	Logger.SetPrefix("INFO: ")
	Logger.Println(message)
}
