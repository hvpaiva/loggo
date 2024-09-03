# Loggo
[![Go Reference](https://pkg.go.dev/badge/github.com/hvpaiva/loggo#section-readme.svg)](https://pkg.go.dev/github.com/hvpaiva/loggo#section-readme)
[![Go Report Card](https://goreportcard.com/badge/github.com/hvpaiva/loggo)](https://goreportcard.com/report/github.com/hvpaiva/loggo)
[![License](https://img.shields.io/badge/License-Mit-blue.svg)](LICENSE)
[![CI](https://github.com/hvpaiva/loggo/actions/workflows/ci.yml/badge.svg)](https://github.com/hvpaiva/loggo/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/hvpaiva/loggo/branch/main/graph/badge.svg)](https://codecov.io/gh/hvpaiva/loggo)


Loggo is a simple and flexible logging library for Go, providing configurable log levels, output destinations, message templates, and time providers.

> **This library has no external dependencies.**

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
- [Thread-Safe Logging](#thread-safe-logging)
- [Vs Go Standard Library Log](#vs-go-standard-library-log)
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

To install Loggo, use `go get`:

```sh
go get github.com/hvpaiva/loggo
```

## Usage

#### Basic Usage

Create a new logger with a specified log level and log messages:

```go
package main

import (
    "os"


    "github.com/hvpaiva/loggo"
)

func main() {
    logger := loggo.New(loggo.LevelInfo)
	
    logger.Info("This is an info message")
    logger.Debug("This debug message will not be logged")
	// Output: 
	// 2024-09-03 15:04:05 [ INFO]: This is an info message
	
}
```

#### Custom Output

Log to a file instead of the standard output:

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

#### Custom Template

Customize the log message format:

```go
package main

import (
    "os"


    "github.com/hvpaiva/loggo"
)

func main() {
    logger := loggo.New(loggo.LevelInfo, loggo.WithTemplate("{{.Time}} - {{.Message}}"))
    logger.Info("This is an info message")
	// Output: 2024-09-03 15:04:05 - This is an info message
}
```

> The template can receive the following placeholders:
> - {{.Level}}: the log level (e.g., "INFO", "DEBUG", etc.)
> - {{.Time}}: the log timestamp (e.g., "2024-09-03 15:04:05")
> - {{.Message}}: the log message
> - {{.Caller}}: the log caller (e.g., "main.go:10")
> 
> The default template is `{{.Time}} [{{printf \"%5s\" .Level}}]: {{.Message}}`.

#### Custom Time Provider

Use a custom time provider for log timestamps:

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

#### Custom Time Format

Use a custom time format for log timestamps:

```go
package main

import (
    "os"


    "github.com/hvpaiva/loggo"
)

func main() {
	logger := loggo.New(
		loggo.LevelInfo,
		loggo.WithTimeFormat("02/01/2006 15:04:05"),
	)
    logger.Info("This message will have a custom time format")
	// Output: 01/02/2006 00:00:00 [ INFO]: This message will have a custom time format
}
```

#### Maximum Log Message Size

Configure the maximum size of a log message:

```go
package main

import (
    "os"


    "github.com/hvpaiva/loggo"
)

func main() {
	logger := loggo.New(
		loggo.LevelInfo,
		loggo.WithMaxSize(10),
	)
    logger.Info("This is an info message with a maximum size")
	// Output: 2024-09-03 15:04:05 [ INFO]: This is an
}
```

#### Custom Caller Provider

Use a custom caller provider for log caller:

```go
package main

import (
    "os"
	
    "github.com/hvpaiva/loggo"
)

func main() {
    logger := loggo.New(
        loggo.LevelInfo,
        loggo.WithCallerProvider(func() (pc uintptr, file string, line int, ok bool) {
            // Implement your custom caller provider here
            return 0, "custom caller", 0, true
        }),
		loggo.WithTemplate("{{.Caller}} - {{.Message}}"),
    )
    logger.Info("This is an info message with a custom caller")
    // Output: custom caller:0 - This is an info message with a custom caller
}
```

> The signature of the caller provider are the same as the [runtime.Caller](https://golang.org/pkg/runtime/#Caller) function.
> - pc uintptr: the program counter -> This value is not used by the caller provider
> - file string: the file name
> - line int: the line number
> - ok bool: a boolean indicating if the information was retrieved successfully
> - skip int: the number of stack frames to skip before getting the caller information

### Thread-Safe Logging

Loggo ensures thread-safe logging using a mutex:

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

### Vs Go Standard Library Log

Loggo provides several advantages over the [Go standard library log](https://pkg.go.dev/log) package:

| Feature                                  | Loggo | Go Standard Library Log |
|------------------------------------------|-------|-------------------------|
| Configurable log levels                  | Yes   | No*                     |
| Customizable output destinations         | Yes   | Yes                     |
| Flexible message templates               | Yes   | No**                    |
| Formatted log messages                   | Yes   | Yes                     |
| Custom time providers for log timestamps | Yes   | No                      |
| Custom time formats for log timestamps   | Yes   | No                      |
| Configurable maximum log message size    | Yes   | No                      |
| Thread-safe logging                      | Yes   | Yes                     |
| Custom caller provider for log caller    | Yes   | No*\**                  |

> \* The Go standard library log package provides a single global logger with no log levels. There are three predefined loggers: 
> `log.Print`, `log.Fatal`, and `log.Panic`, which fell like a log level, but they are more like log actions.
> (Just print, print and exit, print and panic, respectively).
>
> \** The Go standard library log package provides a set of predefined message formats that cannot be customized.
>
> \*\** The Go standard library log package provides caller information in the log message, but it is not configurable.

> 
## Documentation

For more detailed documentation and examples, please refer to the [GoDoc](https://pkg.go.dev/github.com/hvpaiva/loggo).

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
