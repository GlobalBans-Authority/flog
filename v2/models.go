package flog

type Severity int

const (
	SeverityDebug Severity = iota
	SeverityInfo
	SeverityWarn
	SeverityError
	SeverityPanic
)

type LogLevel string

const (
	LogPanic   LogLevel = "panic"
	LogError   LogLevel = "error"
	LogWarn    LogLevel = "warn"
	LogInfo    LogLevel = "info"
	LogDebug   LogLevel = "debug"
	LogSuccess LogLevel = "success"
	LogReader  LogLevel = "reader"
	LogChannel LogLevel = "channel"
)

type Fields map[string]interface{}

type Color string

type RGB struct {
	R int `json:"R"`
	G int `json:"G"`
	B int `json:"B"`
}

type LogLevelConfig struct {
	Color        Color
	LogToConsole bool
	LogToFile    bool
	FileFolder   string
	Severity     Severity
}

type CallerInfo struct {
	funcName string
	line     int
}

type LevelLogger struct {
	logger *Logger
	level  LogLevel
}
