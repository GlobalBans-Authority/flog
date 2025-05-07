package flog

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type Config struct {
	LogFolder     string
	Colors        Colors
	LogConsole    bool
	LogFilePrefix string
	FormatPrefix  string
	MinSeverity   Severity // Minimum severity to log
	LogFormat     string   // Custom log format
	MaxLogSize    int64    // Max file size in bytes before rotation
}

type Colors struct {
	LogPanic   Color `json:"log_panic,omitempty"`
	LogError   Color `json:"log_error,omitempty"`
	LogWarn    Color `json:"log_warn,omitempty"`
	LogInfo    Color `json:"log_info,omitempty"`
	LogSuccess Color `json:"log_success,omitempty"`
	LogDebug   Color `json:"log_debug,omitempty"`
}

func (co Colors) Default() Colors {
	co.LogPanic = AnsiRGB(RGB{R: 255, G: 0, B: 0})
	co.LogError = AnsiRGB(RGB{R: 234, G: 1, B: 1})
	co.LogWarn = AnsiRGB(RGB{R: 234, G: 173, B: 1})
	co.LogInfo = AnsiRGB(RGB{R: 0, G: 86, B: 234})
	co.LogSuccess = AnsiRGB(RGB{R: 1, G: 235, B: 110})
	co.LogDebug = AnsiRGB(RGB{R: 128, G: 128, B: 128})
	return co
}

func Default() Config {
	c := Config{}
	folder, _ := os.UserCacheDir()
	logFolder := filepath.Join(folder, "FLog")
	c.LogFolder = path.Join(logFolder, "logs")
	c.LogConsole = true
	c.Colors = Colors{}.Default()
	c.FormatPrefix = "!"
	c.LogFilePrefix = ""
	c.MinSeverity = SeverityInfo
	c.LogFormat = "[ {timestamp} ] [ {caller_func} → {caller_line} ]: {message} {fields}" // Restore original format
	c.MaxLogSize = 10 * 1024 * 1024                                                       // 10MB
	return c
}

func AnsiRGB(rgb RGB) Color {
	return Color("\x1b[38;2;" + strconv.FormatInt(int64(rgb.R), 10) + ";" + strconv.FormatInt(int64(rgb.G), 10) + ";" + strconv.FormatInt(int64(rgb.B), 10) + "m")
}

const resetColor = "\x1b[0m"
