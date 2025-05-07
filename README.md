# FLog - A Simple Farmers Logger
<p align="center"><img src = "/assets/FLog.png"></p>

<p align="center"><img src = "/assets/example.png"></p>

FLog is a straightforward logger for Go that keeps things simple, just like logging should be. This is a fork of FLog, The GlobalBans logger, full up to date one can be found at https://github.com/GlobalBans-Authority/FLog

### Features:
- Multiple log levels (Error, Warn, Info, Debug, Success)
- Color-coded console output
- Concurrent file logging with buffered writes
- Automatic log file rotation
- Silent logging option
- Caller information tracking
- Formatted logging support
- Built to "Make Loggers Great Again."



## Upcoming features
- Log rotation
- Log Sanitization
- Level filtering
- Configs for all above


## Usage

**GO Version:** 1.22 windows/amd64

```diff
! Right now this is a prototype, a proof of concept. i would not advise using this in production!
```
seriously its not ready lol


## Quick Start

### Installing

```go
go get github.com/Make-Tarkov-Great-Again/flog/v4/flog@latest
```
### Basic Usage

```go
// Basic logging
import "github.com/Make-Tarkov-Great-Again/flog/v4/flog"

...
flog.Init(flog.Default()) // Initialize the logger with default settings (Fastest)
flog.Info("Starting application...")
flog.Error("An error occurred:", err)
```

## Known issues

Working on it...

## Who are we, and why are we doing this?

Make Tarkov Great Again is a project founded by a few Altered Escape developers. We felt that there shouldn't just be one option for Tarkov offline play and began development on MTGA.


## Contribution

- Is there a part of FLog you want to tackle?
- Some code you would like to refactor?
- Got an idea you would like to share/implement?

Feel free to create a fork, open a pull request, and request a review: **We are open to any contribution!**

**To keep your fork up-to-date, we recommend using Pull!**

<p align="center"><img src = "https://user-images.githubusercontent.com/21200584/183050357-6c92f1cd-68ca-4f74-b41d-1706915c67cf.gif"></p>
