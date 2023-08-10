package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const logsFolder = "Logs"
const loggingRoutesFile = "Logs/Logging-Routes.gbconfig"

var logFileMap = make(map[string]string)
var flog = make(map[string]func(message string))

func GetFlog() map[string]func(message string) {
	return flog
}
func logDating() string {
	now := time.Now()
	milliseconds := now.Nanosecond() / 1e6
	return fmt.Sprintf("[%02d:%02d:%02d:%03d] ->", now.Hour(), now.Minute(), now.Second(), milliseconds)
}

func updateLoggingRoutesFile() {
	logTypeFolders := make(map[string]string)

	// Generate log type folders for each log type
	logTypes := []string{"warn", "error", "info", "debug", "security", "network", "incident", "discord", "silent"}
	for _, logType := range logTypes {
		logTypeFolder := path.Join(logsFolder, logType)
		logTypeFolders[logType] = logTypeFolder

		if _, err := os.Stat(logTypeFolder); os.IsNotExist(err) {
			os.Mkdir(logTypeFolder, 0755)
		}
	}

	// Update the log file map with new log file paths
	for logType, logTypeFolder := range logTypeFolders {
		if _, exists := logFileMap[logType]; !exists {
			logFileName := fmt.Sprintf("log_%s_%d.log", logType, time.Now().UnixNano())
			logFilePath := path.Join(logTypeFolder, logFileName)
			logFileMap[logType] = logFilePath
		}
	}

	// Update the logging routes file
	var routes []string
	for logType, logFileName := range logFileMap {
		routes = append(routes, fmt.Sprintf("%s:%s", logType, logFileName))
	}
	routesContent := strings.Join(routes, "\n")
	err := os.WriteFile(loggingRoutesFile, []byte(routesContent), 0644)
	if err != nil {
		log.Printf("Failed to update logging routes file: %v\n", err)
	}
}

func readLoggingRoutesFile() {
	content, err := os.ReadFile(loggingRoutesFile)
	if err != nil {
		log.Printf("Failed to read logging routes file: %v\n", err)
		return
	}

	routes := strings.Split(string(content), "\n")
	for _, route := range routes {
		parts := strings.Split(route, ":")
		if len(parts) == 2 {
			logFileMap[parts[0]] = parts[1]
		}
	}
}

func createLogFunction(logType string) func(message string) {
	readLoggingRoutesFile()
	return func(message string) {
		flogg(logType, message)
	}
}

func flogg(logType, message string) {
	if _, exists := logFileMap[logType]; !exists {
		logTypeFolder := path.Join(logsFolder, logType)
		if _, err := os.Stat(logTypeFolder); os.IsNotExist(err) {
			os.Mkdir(logTypeFolder, 0755)
		}

		logFileName := fmt.Sprintf("log_%s_%d.log", logType, time.Now().UnixNano())
		logFilePath := path.Join(logTypeFolder, logFileName)
		logFileMap[logType] = logFilePath
	}

	logFile, err := os.OpenFile(logFileMap[logType], os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	prefix := fmt.Sprintf("%s [%s]:", logDating(), strings.ToUpper(logType))
	log.Println(prefix, message)
}

func init() {
	updateLoggingRoutesFile()
	readLoggingRoutesFile()

	flog = map[string]func(message string){
		"warn":     createLogFunction("warn"),
		"error":    createLogFunction("error"),
		"info":     createLogFunction("info"),
		"debug":    createLogFunction("debug"),
		"security": createLogFunction("security"),
		"network":  createLogFunction("network"),
		"incident": createLogFunction("incident"),
		"discord":  createLogFunction("discord"),
		"silent":   createLogFunction("silent"),
	}
}
