# FLog - A Simple Farmers Logger
<p align="center"><img src = "/assets/FLog.png"></p>

<p align="center"><img src = "/assets/example.png"></p>

FLog is a straightforward logger for Go that keeps things simple, just like logging should be.

### Features:
- Multiple log levels (Error, Warn, Info, Debug, Success)
- Color-coded console output from RGB
- Concurrent file logging with buffered writes
- Automatic log file rotation
- Silent logging option
- Caller information tracking
- Formatted logging support
- Built to "Make Loggers Great Again."
- Level Sanitization
- Custom log levels
- yada yada no one reads this far into the read me anyways


most this above should work, what do i look like to you? a devoloper?


## Upcoming features
- Log Sanitization
- Configs for all above


## Usage

**GO Version:** 1.22 windows/amd64




## Quick Start

### Installing

```go
go get github.com/GlobalBans-Authority/flog/v2@latest
```
### Basic Usage

```go
// Basic logging
import "github.com/GlobalBans-Authority/flog/v2"

...
flog.Init(flog.Default()) // Initialize the logger with default settings (Fastest)
flog.Info("Starting application...")
flog.Error("An error occurred:", err)
```

## Known issues

Working on it...


## Contribution

- Is there a part of FLog you want to tackle?
- Some code you would like to refactor?
- Got an idea you would like to share/implement?

Feel free to create a fork, open a pull request, and request a review: **We are open to any contribution!**

**To keep your fork up-to-date, we recommend using Pull!**