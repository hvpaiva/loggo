# Loggo - A Go Logging Library

[![Go Reference](https://pkg.go.dev/badge/github.com/hvpaiva/loggo#section-readme.svg)](https://pkg.go.dev/github.com/hvpaiva/loggo#section-readme)
[![License](https://img.shields.io/badge/License-Mit-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hvpaiva/loggo)](https://goreportcard.com/report/github.com/hvpaiva/loggo)
[![codecov](https://codecov.io/gh/hvpaiva/loggo/branch/main/graph/badge.svg)](https://codecov.io/gh/hvpaiva/loggo)

[![CI](https://github.com/hvpaiva/loggo/actions/workflows/ci.yml/badge.svg)](https://github.com/hvpaiva/loggo/actions/workflows/ci.yml)
[![CodeQL](https://github.com/hvpaiva/loggo/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/hvpaiva/loggo/actions/workflows/github-code-scanning/codeql)
[![Dependabot Updates](https://github.com/hvpaiva/loggo/actions/workflows/dependabot/dependabot-updates/badge.svg)](https://github.com/hvpaiva/loggo/actions/workflows/dependabot/dependabot-updates)

```
 __         ______     ______     ______     ______    
/\ \       /\  __ \   /\  ___\   /\  ___\   /\  __ \   
\ \ \____  \ \ \/\ \  \ \ \__ \  \ \ \__ \  \ \ \/\ \  
 \ \_____\  \ \_____\  \ \_____\  \ \_____\  \ \_____\ 
  \/_____/   \/_____/   \/_____/   \/_____/   \/_____/ 
                                                       
```

Loggo is a simple and flexible logging library for Go, offering configurable log levels, output destinations, message 
templates, time providers, and more.

> **Note:** This library has no external dependencies.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Basic Usage](#basic-usage)
  - [Custom Output](#custom-output)
  - [Custom Template](#custom-template)
  - [Custom Time Provider](#custom-time-provider)
  - [Custom Time Format](#custom-time-format)
  - [Maximum Log Message Size](#maximum-log-message-size)
  - [Custom Caller Provider](#custom-caller-provider)
  - [Context Logging](#context-logging)
  - [Pre & Post Log Hooks](#pre--post-log-hooks)
- [Thread-Safe Logging](#thread-safe-logging)
- [Comparison with Go's Standard Library](#comparison-with-gos-standard-library)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## Features

- Configurable log levels (Debug, Info, Warn, Error, Fatal)
- Customizable output destinations (e.g., `os.Stdout`, `os.Stderr`, files)
- Flexible message templates
- Custom time providers for log timestamps
- Custom time formats for log timestamps
- Configurable maximum log message size
- Thread-safe logging

## Installation

To install Loggo, run:

```sh
go get github.com/hvpaiva/loggo
```

## Usage

### Basic Usage

Initialize a logger with a specified log level and log messages:

```go
package main

import "github.com/hvpaiva/loggo"

func main() {
    logger := loggo.New(loggo.LevelInfo)
	
    logger.Info("This is an info message")
    logger.Debug("This debug message will not be logged")
    // Output: 2024-09-03 15:04:05 [ INFO]: This is an info message
}
```

### Custom Output

Redirect logs to a file instead of standard output:

```go
package main

import (
    "os"


    "github.com/hvpaiva/loggo"
)

func main() {
    file, _ := os.Create("log.txt")
    logger := loggo.New(loggo.LevelInfo, loggo.WithOutput(file))
	
    logger.Info("This message will be logged to a file")
	// In the file log.txt:
	// Output: 2024-09-03 15:04:05 [ INFO]: This message will be logged to a file
}
```

### Custom Template

Define a custom format for log messages:

```go
package main

import "github.com/hvpaiva/loggo"

func main() {
    logger := loggo.New(loggo.LevelInfo, loggo.WithTemplate("{{.Time}} - {{.Message}}"))
    logger.Info("This is an info message")
    // Output: 2024-09-03 15:04:05 - This is an info message
}
```

> Available template placeholders:
> - `{{.Level}}`: log level (e.g., "INFO", "DEBUG")
> - `{{.Time}}`: log timestamp (e.g., "2024-09-03 15:04:05")
> - `{{.Message}}`: log message
> - `{{.Caller}}`: log caller (e.g., "main.go:10")
>
> Default template: `{{.Time}} [{{printf \"%5s\" .Level}}]: {{.Message}}`.

### Custom Time Provider

Specify a custom time provider for timestamps:

```go
package main

import (
    "time"


    "github.com/hvpaiva/loggo"
)

func fakeNow() time.Time {
    return time.Unix(0, 0)
}

func main() {
    logger := loggo.New(loggo.LevelInfo, loggo.WithTimeProvider(fakeNow))
    logger.Info("This message will have a custom timestamp")
    // Output: 1970-01-01 00:00:00 [ INFO]: This message will have a custom timestamp
}
```

### Custom Time Format

Set a custom time format for timestamps:

```go
package main

import "github.com/hvpaiva/loggo"

func main() {
	logger := loggo.New(
		loggo.LevelInfo,
		loggo.WithTimeFormat("02/01/2006 15:04:05"),
	)
    logger.Info("This message will have a custom time format")
    // Output: 01/02/2006 15:04:05 [ INFO]: This message will have a custom time format
}
```

### Maximum Log Message Size

Limit the maximum size of a log message:

```go
package main

import "github.com/hvpaiva/loggo"

func main() {
	logger := loggo.New(
		loggo.LevelInfo,
		loggo.WithMaxSize(10),
	)
    logger.Info("This is an info message with a maximum size")
    // Output: 2024-09-03 15:04:05 [ INFO]: This is an
}
```

### Custom Caller Provider

Use a custom caller provider for log caller information:

```go
package main

import "github.com/hvpaiva/loggo"

func main() {
    logger := loggo.New(
        loggo.LevelInfo,
        loggo.WithCallerProvider(func() (pc uintptr, file string, line int, ok bool) {
            return 0, "custom caller", 0, true
        }),
		loggo.WithTemplate("{{.Caller}} - {{.Message}}"),
    )
    logger.Info("This is an info message with a custom caller")
    // Output: custom caller:0 - This is an info message with a custom caller
}
```

### Context Logging

Enhance a logger with additional context using `WithContext`:

```go
package main

import (
    "context"

    "github.com/hvpaiva/loggo"
)

func main() {
    ctx := context.WithValue(context.Background(), "trace_id", "123456")

    logger := loggo.New(loggo.LevelInfo, loggo.WithContext(ctx))
    logger.Infof("This is an info message with context, trace_id: %s", logger.Context.Value("trace_id"))
    // Output: 2024-09-03 15:04:05 [ INFO]: This is an info message with context, trace_id: 123456
}
```

### Pre & Post Log Hooks

Execute custom logic before and after a log message:

```go
package main

import (
    "fmt"
    "strings"

    "github.com/hvpaiva/loggo"
)

func main() {
    logger := loggo.New(
        loggo.LevelInfo,
        loggo.WithPreHook(func(l *loggo.Logger, msg *string) {
            *msg = strings.ToUpper(*msg)
        }),
        loggo.WithPostHook(func(l *loggo.Logger, msg *string) {
            fmt.Println("Log message written:", *msg)
        }),
    )

    logger.Info("This is an info message")
    // Output: 2024-09-03 15:04:05 [ INFO]: THIS IS AN INFO MESSAGE
    // Log message written: THIS IS AN INFO MESSAGE
}
```

## Thread-Safe Logging

Loggo ensures thread safety using a mutex:

```go
package main

import (
  "os"
  "sync"

  "github.com/hvpaiva/loggo"
)

func main() {
  logger := loggo.New(loggo.LevelInfo, loggo.WithOutput(os.Stdout))
  var wg sync.WaitGroup

  for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(i int) {
      defer wg.Done()
      logger.Infof("Logging from goroutine %d", i)
    }(i)
  }

  wg.Wait()
}
```

## Comparison with Go's Standard Library

Loggo provides several advantages over the [Go standard library log package](https://pkg.go.dev/log):

| Feature                                  | Loggo | Go Log |
|------------------------------------------|-------|--------|
| Configurable log levels                  | Yes   | No*    |
| Customizable output destinations         | Yes   | Yes    |
| Flexible message templates               | Yes   | No**   |
| Formatted log messages                   | Yes   | Yes    |
| Custom time providers                    | Yes   | No     |
| Custom time formats                      | Yes   | No     |
| Maximum log message size                 | Yes   | No     |
| Thread-safe logging                      | Yes   | Yes    |
| Custom caller provider                   | Yes   | No*\** |

> \* Go's log level is action-based with functions like `log.Print`, `log.Fatal`, and `log.Panic`.
> 
> \** Predefined message format in Go's log.
> 
> \*\** Caller info in Go's log is not configurable.

## Documentation

For detailed documentation and examples, visit the [GoDoc](https://pkg.go.dev/github.com/hvpaiva/loggo).

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
