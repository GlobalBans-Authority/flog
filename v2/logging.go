package flog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	Config       Config
	logFolder    string
	logFileMap   map[string]*bufio.Writer
	files        map[string]*os.File
	mu           sync.Mutex
	bufPool      sync.Pool
	levelConfigs map[LogLevel]LogLevelConfig
	fields       Fields
	stopChan     chan struct{}
	watchWg      sync.WaitGroup
}

var logger *Logger

func Init(config Config) {
	logger = &Logger{
		logFolder:    path.Join(config.LogFolder, "logs"),
		logFileMap:   make(map[string]*bufio.Writer),
		files:        make(map[string]*os.File),
		levelConfigs: make(map[LogLevel]LogLevelConfig),
		fields:       make(Fields),
		stopChan:     make(chan struct{}),
		bufPool: sync.Pool{
			New: func() interface{} {
				return new(strings.Builder)
			},
		},
	}
	if config.Colors == (Colors{}) {
		config.Colors = Colors{}.Default()
	}
	logger.Config = config
	logger.levelConfigs = map[LogLevel]LogLevelConfig{
		LogPanic:   {Color: config.Colors.LogPanic, LogToConsole: true, LogToFile: true, FileFolder: "error", Severity: SeverityPanic},
		LogError:   {Color: config.Colors.LogError, LogToConsole: true, LogToFile: true, FileFolder: "error", Severity: SeverityError},
		LogWarn:    {Color: config.Colors.LogWarn, LogToConsole: true, LogToFile: true, FileFolder: "warn", Severity: SeverityWarn},
		LogInfo:    {Color: config.Colors.LogInfo, LogToConsole: true, LogToFile: true, FileFolder: "info", Severity: SeverityInfo},
		LogDebug:   {Color: config.Colors.LogDebug, LogToConsole: true, LogToFile: true, FileFolder: "debug", Severity: SeverityDebug},
		LogSuccess: {Color: config.Colors.LogSuccess, LogToConsole: true, LogToFile: true, FileFolder: "info", Severity: SeverityInfo},
		LogReader:  {Color: config.Colors.LogInfo, LogToConsole: true, LogToFile: true, FileFolder: "reader", Severity: SeverityInfo},
		LogChannel: {Color: config.Colors.LogInfo, LogToConsole: true, LogToFile: true, FileFolder: "channel", Severity: SeverityInfo},
	}
	_ = logger.initFolders()
	_ = logger.initLogFiles()
	go logger.periodicFlush(logger.stopChan)
}

func (l *Logger) WithFields(fields Fields) *Logger {
	newFields := make(Fields)
	for k, v := range l.fields {
		newFields[k] = v
	}
	for k, v := range fields {
		newFields[k] = v
	}
	return &Logger{
		Config:       l.Config,
		logFolder:    l.logFolder,
		logFileMap:   l.logFileMap,
		files:        l.files,
		mu:           l.mu,
		bufPool:      l.bufPool,
		levelConfigs: l.levelConfigs,
		fields:       newFields,
		stopChan:     l.stopChan,
	}
}

func (l *Logger) RegisterLogLevel(level LogLevel, config LogLevelConfig) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Validate custom level
	if _, exists := l.levelConfigs[level]; exists && level != LogPanic && level != LogError && level != LogWarn && level != LogInfo && level != LogDebug && level != LogSuccess {
		fmt.Fprintf(os.Stderr, "Warning: Overwriting existing log level %s\n", level)
	}

	l.levelConfigs[level] = config

	if config.LogToFile {
		if _, exists := l.logFileMap[config.FileFolder]; !exists {
			logTypeFolder := path.Join(l.logFolder, config.FileFolder)
			if _, err := os.Stat(logTypeFolder); os.IsNotExist(err) {
				if err := os.MkdirAll(logTypeFolder, 0755); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to create folder %s: %v\n", logTypeFolder, err)
					return
				}
			}

			logFileName := fmt.Sprintf("log_%s_%d.%s", config.FileFolder, time.Now().UnixNano(), l.Config.LogFilePrefix)
			logFilePath := path.Join(logTypeFolder, logFileName)
			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open file %s: %v\n", logFilePath, err)
				return
			}

			l.files[config.FileFolder] = file
			l.logFileMap[config.FileFolder] = bufio.NewWriterSize(file, 32*1024)
		}
	}
}

func (l *Logger) GetLogger(level LogLevel) *LevelLogger {
	return &LevelLogger{logger: l, level: level}
}

func (l *Logger) log(logType LogLevel, format string, args ...interface{}) {
	config, exists := l.levelConfigs[logType]
	if !exists || config.Severity < l.Config.MinSeverity {
		return
	}

	silent := false
	if len(args) > 0 {
		if val, ok := args[len(args)-1].(bool); ok {
			silent = val
			args = args[:len(args)-1]
		}
	}
	formattedMsg := fmt.Sprintf(format, args...)
	logEntry := l.formatLogEntry(logType, formattedMsg, l.fields)

	if config.LogToFile {
		l.mu.Lock()
		if writer, ok := l.logFileMap[config.FileFolder]; ok {
			if _, err := writer.WriteString(logEntry); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write to %s: %v\n", config.FileFolder, err)
			}
			// Immediate flush to ensure logs are written
			if err := writer.Flush(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to flush %s: %v\n", config.FileFolder, err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "No writer found for folder %s\n", config.FileFolder)
		}
		l.mu.Unlock()
	}

	if config.LogToConsole && !silent && l.Config.LogConsole {
		prefix := fmt.Sprintf("[%s]", strings.ToUpper(string(logType)))
		st := fmt.Sprintf("%s%s %s%s", config.Color, prefix, resetColor, logEntry)
		os.Stdout.WriteString(st)
	}
}

func (ll *LevelLogger) Log(format string, args ...interface{}) {
	ll.logger.log(ll.level, format, args...)
}

func (l *Logger) formatLogEntry(logType LogLevel, message string, fields Fields) string {
	format := l.Config.LogFormat
	if format == "" {
		format = "[{timestamp}] [{caller_func} → {caller_line}]: {message} {fields}"
	}
	timestamp := l.logDating()
	caller := getCallerInfo(4)
	fieldsStr := ""
	if fields != nil {
		var sb strings.Builder
		for k, v := range fields {
			if sb.Len() > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(fmt.Sprint(v))
		}
		fieldsStr = sb.String()
	}
	entry := strings.ReplaceAll(format, "{timestamp}", timestamp)
	entry = strings.ReplaceAll(entry, "{level}", string(logType))               // Use entry instead of format
	entry = strings.ReplaceAll(entry, "{message}", message)                     // Use entry instead of format
	entry = strings.ReplaceAll(entry, "{fields}", fieldsStr)                    // Use entry instead of format
	entry = strings.ReplaceAll(entry, "{caller_func}", caller.funcName)         // Use entry instead of format
	entry = strings.ReplaceAll(entry, "{caller_line}", fmt.Sprint(caller.line)) // Use entry instead of format
	return entry + "\n"
}

func (l *Logger) periodicFlush(stopChan chan struct{}) {
	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-ticker.C:
			l.mu.Lock()
			for folder, writer := range l.logFileMap {
				if err := writer.Flush(); err != nil {
					fmt.Fprintf(os.Stderr, "Flush error for %s: %v\n", folder, err)
				}
				file := l.files[folder]
				if stat, err := file.Stat(); err == nil && stat.Size() > l.Config.MaxLogSize {
					file.Close()
					newFileName := fmt.Sprintf("log_%s_%d.%s", folder, time.Now().UnixNano(), l.Config.LogFilePrefix)
					newFilePath := path.Join(l.logFolder, folder, newFileName)
					newFile, err := os.OpenFile(newFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Rotation error for %s: %v\n", folder, err)
						continue
					}
					l.files[folder] = newFile
					l.logFileMap[folder] = bufio.NewWriterSize(newFile, 32*1024)
				}
			}
			l.mu.Unlock()
		case <-stopChan:
			ticker.Stop()
			return
		}
	}
}

func (l *Logger) Cleanup() {
	l.stopChan <- struct{}{}
	l.watchWg.Wait() // Wait for all watchers to finish
	l.mu.Lock()
	defer l.mu.Unlock()
	for folder, writer := range l.logFileMap {
		if err := writer.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "Cleanup flush error for %s: %v\n", folder, err)
		}
	}
	for folder, file := range l.files {
		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Cleanup close error for %s: %v\n", folder, err)
		}
	}
}

func Cleanup() {
	if logger != nil {
		logger.Cleanup()
	}
}

func (l *Logger) initFolders() error {
	if _, err := os.Stat(l.logFolder); os.IsNotExist(err) {
		return os.MkdirAll(l.logFolder, 0755)
	}
	return nil
}

func (l *Logger) initLogFiles() error {
	l.logFileMap = make(map[string]*bufio.Writer)
	l.files = make(map[string]*os.File)

	folders := make(map[string]bool)
	for _, cfg := range l.levelConfigs {
		if cfg.LogToFile {
			folders[cfg.FileFolder] = true
		}
	}

	for folder := range folders {
		logTypeFolder := path.Join(l.logFolder, folder)
		if _, err := os.Stat(logTypeFolder); os.IsNotExist(err) {
			if err := os.MkdirAll(logTypeFolder, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create folder %s: %v\n", logTypeFolder, err)
				return err
			}
		}
		fmt.Fprintf(os.Stderr, "Initialized folder: %s\n", logTypeFolder)
		logFileName := fmt.Sprintf("log_%s_%d.%s", folder, time.Now().UnixNano(), l.Config.LogFilePrefix)
		logFilePath := path.Join(logTypeFolder, logFileName)
		fmt.Fprintf(os.Stderr, "Initialized log file: %s\n", logFilePath)

		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open file %s: %v\n", logFilePath, err)
			return err
		}

		l.files[folder] = file
		l.logFileMap[folder] = bufio.NewWriterSize(file, 32*1024)
	}

	return nil
}

var timeFormat = "15:04:05.000"

func (l *Logger) logDating() string {
	return time.Now().Format(timeFormat)
}

var (
	callerCache = make(map[uintptr]CallerInfo)
	callerMu    sync.RWMutex
)

func getCallerInfo(skip int) CallerInfo {
	//defer measurment.Un(measurment.Trace("Get caller info"))

	pc, _, line, _ := runtime.Caller(skip)

	callerMu.RLock()
	if info, ok := callerCache[pc]; ok {
		callerMu.RUnlock()
		info.line = line // Update line number
		return info
	}
	callerMu.RUnlock()

	callerFunc := runtime.FuncForPC(pc)
	if callerFunc == nil {
		return CallerInfo{"Anonymous", line}
	}

	fullName := callerFunc.Name()

	callerMu.Lock()
	callerCache[pc] = CallerInfo{fullName, line}
	info := callerCache[pc]
	callerMu.Unlock()

	return info
}

func GetFlogger() *Logger {
	return logger
}

// WatchReader reads from an io.Reader and logs each line
func (l *Logger) WatchReader(r io.Reader, level LogLevel) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		l.log(level, "%s", scanner.Text())
	}
}

// WatchChannel listens to a channel and logs each message
func (l *Logger) WatchChannel(ch interface{}, level LogLevel) {
	val := reflect.ValueOf(ch)
	if val.Kind() != reflect.Chan {
		l.Error("WatchChannel requires a channel")
		return
	}

	l.watchWg.Add(1)
	go func() {
		defer l.watchWg.Done()
		for {
			chosen, recv, ok := reflect.Select([]reflect.SelectCase{
				{
					Dir:  reflect.SelectRecv,
					Chan: val,
				},
				{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(l.stopChan),
				},
			})

			if chosen == 1 || !ok { // stopChan received or channel closed
				return
			}

			l.log(level, fmt.Sprintf("%s", recv.String()))
			fmt.Println("logging it", fmt.Sprintf("%s", recv.String()))
		}
	}()
}
