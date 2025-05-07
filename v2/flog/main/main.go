package main

import (
	"fmt"
	"os"
	"time"

	flog "github.com/GlobalBans-Authority/v2/flog"
)

func main() {
	wd, _ := os.Getwd()
	folder := wd
	fmt.Println("Logging folder ", folder)

	flog.Init(flog.Config{
		LogFilePrefix: "ballin",
		LogFolder:     folder,
		LogConsole:    true,
		FormatPrefix:  "%",
		MinSeverity:   1,
	})
	logger := flog.GetFlogger()

	logger.RegisterLogLevel("trace", flog.LogLevelConfig{
		Color:        flog.AnsiRGB(flog.RGB{R: 100, G: 100, B: 100}),
		LogToConsole: true,
		LogToFile:    true,
		Severity:     1,
		FileFolder:   "trace",
	})

	ch := make(chan string, 10)
	flog.WatchChannel(ch, flog.LogInfo)

	ch <- "Test message 1"
	ch <- "Test message 2"

	for i := 0; i < 50; i++ {
		ch <- "test"
	}
	time.Sleep(time.Second * 5) //give some time for messages to be processed
	close(ch)
}
