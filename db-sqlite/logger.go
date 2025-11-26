package main

import (
	"fmt"
	"log"
	"os"
)

// WriteAppLog writes a log message to app.log and optionally prints to stdout
func WriteAppLog(str string, print bool) {
	logPath := "app.log"
	logfile, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		if print {
			fmt.Printf("Failed to open log file: %v\n", err)
		}
		return
	}
	defer logfile.Close()

	// Create a new logger to write to the file
	fileLogger := log.New(logfile, "", log.LstdFlags)
	fileLogger.Println(str)

	if print {
		fmt.Println(str)
	}
}

