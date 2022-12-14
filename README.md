# Common Go Logger

Standard logging library for golang projects

## Installation

In order to use the library, you must create the configuration, a new logger instance and the methods you want to use.

More information about the methods and their uses below.

Example:

```go
package logger

import (
	"context"
	"net/http"
	"os"

	logging "github.com/LemontechSA/common-go-logger"
)

var log logging.Logger

func init() {
  configuration := logging.Configuration{
    Environment:  "development",
    Service:      "plutarch",
    Team:         "surveycorps",
    Project:      "te2eus",
    ConsoleLevel: "info",
  }

  log = logging.NewLogger(configuration)
}

func Info(description string, message string, payload map[string]string) {
	log.Info(description, message, payload)
}

func Warn(description string, message string, payload map[string]string) {
	log.Warn(description, message, payload)
}

func Debug(description string, message string, payload map[string]string) {
	log.Debug(description, message, payload)
}

func Error(description string, message string, payload map[string]string) {
	log.Error(description, message, payload)
}

func Fatal(description string, message string, payload map[string]string) {
	log.Fatal(description, message, payload)
}

func CreateRequestContext(req *http.Request) context.Context {
	return logging.CreateRequestContext(req)
}

func SetContext(ctx context.Context) {
	log.SetContext(ctx)
}

func ClearContext() {
	log.ClearContext()
}
```

You can use `os.Getenv("ENV_NAME")` to get the configuration value dinamically from the environment vars.

**NOTE: Do not forget to execute the command `go mod tidy` to install the library in the project**

### Configuration Glossary

#### Environment

Environment associated with the event.

#### Service

Name of the service associated with the event.

#### Team

Name of the team in charge of the service.

#### Project

Name of the project to which the service belongs.

#### ConsoleLevel

Level from which the logs will be displayed in the console.

Options:

- info
- warn
- debug
- error
- fatal

The order is descending, so if you put the debug option, they will be displayed from the same level downwards, therefore info and warn will not be displayed.

#### General

If any or all of the options is left empty it will not appear in the log and the default value for `ConsoleLevel` is `info`.

## Usage examples

For all cases, the following fields will be automatically added to the log:

- method: from where the log was called.
- pid: proccess identifier.
- host: host name.

### Basic usage

This applies to all log types, the only difference is the log level.

Example with all the data:

```go
logger.Info("Starting server", "We are ready to GO!", map[string]string{
  "extra1": "test 1",
  "extra2": "test 2",
})
```

The log will be:

```json
{
  "level": "info",
  "date": "2022-12-13T17:16:24.069Z",
  "method": "server/server.go:23",
  "message": "We are ready to GO!",
  "pid": 0,
  "host": "8cfa08dc7116",
  "service": "plutarch",
  "environment": "development",
  "team": "surveycorps",
  "project": "te2eus",
  "description": "Starting server",
  "payload": { "extra1": "test 1", "extra2": "test 2" }
}
```

Example without data:

```go
logger.Info("", "", nil)
```

The log will be:

```json
{
  "level": "info",
  "date": "2022-12-13T17:16:24.069Z",
  "method": "server/server.go:23",
  "message": "",
  "pid": 0,
  "host": "8cfa08dc7116",
  "service": "plutarch",
  "environment": "development",
  "team": "surveycorps",
  "project": "te2eus",
  "description": ""
}
```

The payload data is the only one that will be omitted if it is empty, the description and the message will appear empty.

Additionally, the payload has a function that parses certain fields, currently they are:

- duration

Example:

If the duration field is sent with a string of numbers, it will be parsed to an int.

**NOTE: The duration field must be a string of only numbers, example: "9876", if it contains letters like "9876 ms" the string can't be parsed and the field will be returned with the value 0, this is because the duration field is used to measure and calculate response times and must be an integer**

If can parse it:

```go
logger.Info("", "", map[string]string{
  "duration": "9876"
})
```

The log will be:

```json
{
  "level": "info",
  "date": "2022-12-13T17:16:24.069Z",
  "method": "server/server.go:23",
  "message": "",
  "pid": 0,
  "host": "8cfa08dc7116",
  "service": "plutarch",
  "environment": "development",
  "team": "surveycorps",
  "project": "te2eus",
  "description": "",
  "payload": { "duration": 9876 }
}
```

If can't parse it:

```go
logger.Info("", "", map[string]string{
  "duration": "9876 ms"
})
```

The log will be:

```json
{
  "level": "info",
  "date": "2022-12-13T17:16:24.069Z",
  "method": "server/server.go:23",
  "message": "",
  "pid": 0,
  "host": "8cfa08dc7116",
  "service": "plutarch",
  "environment": "development",
  "team": "surveycorps",
  "project": "te2eus",
  "description": "",
  "payload": { "duration": 0 }
}
```

### With traceability usage

If you want to add the header data in each request, you must use the `CreateRequestContext`, `SetContext` and `ClearContext` methods, creating a middleware that uses these functions so that they are automatically added in each request.

Allowed headers:

- trace_id
- request_id
- session_id
- consumer_name

Any other headers will be ignored.

Example with gin framework:

```go
func LoggerMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx := logger.CreateRequestContext(c.Request)

    logger.SetContext(ctx)

    c.Next()

    logger.ClearContext()
  }
}
```

Router

```go
router.Use(LoggerMiddleware())
...endpoints
```

Example with fiber framework:

```go
func loggerMiddleware(c *fiber.Ctx) error {
	req, _ := http.NewRequest("", "", nil)
	req.Header.Set("trace_id", c.Get("trace_id"))
	req.Header.Set("request_id", c.Get("request_id"))
	req.Header.Set("session_id", c.Get("session_id"))
	req.Header.Set("consumer_name", c.Get("consumer_name"))

	ctx := logger.CreateRequestContext(req)

	logger.SetContext(ctx)

	c.Next()

	logger.ClearContext()

	return nil
}
```

Router

```go
router.Use(LoggerMiddleware)
...endpoints
```

Assuming that we send all the headers in the request headers with some data `test`, when using any log that is within that same request, they will be automatically added to the log.

Example:

```go
logger.Error("Parsing schema", err.Error(), map[string]string{"status": "500"})
```

The log will be:

```json
{
  "level": "error",
  "date": "2022-12-13T17:49:43.829Z",
  "method": "server/bulkload.go:29",
  "message": "json: cannot unmarshal string into Go struct field of type int",
  "pid": 0,
  "host": "8cfa08dc7116",
  "service": "plutarch",
  "environment": "development",
  "team": "surveycorps",
  "project": "te2eus",
  "trace_id": "test",
  "request_id": "test",
  "session_id": "test",
  "consumer_name": "test",
  "description": "Parsing schema",
  "payload": { "status": "500" }
}
```

**IMPORTANT: Do not forget to clean the context at the end of the request, because if it is not done, all subsequent logs will have the values of the previous context.**

**NOTE: Because each framework operates in a different way, the functions were left so that they can be applied in each of them and not a general middleware.**
