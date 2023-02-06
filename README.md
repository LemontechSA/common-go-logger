# Common Go Logger

Standard logging library for golang projects.

## Installation

In order to use the library, you must create the configuration, the loggers instances and the methods you want to use.

If you use the Context logger you must import/export the logging contextKey you want to use.

More information about the loggers, methods and their uses below.

Example:

```go
package logger

import (
	"context"
	"os"

	logging "github.com/LemontechSA/common-go-logger"
)

var genericLog logging.Logger
var contextLog logging.ContextLogger

var ContextKeyCorrelationID = logging.ContextKeyCorrelationID
var ContextKeyCausationID = logging.ContextKeyCausationID
var ContextKeyTenant = logging.ContextKeyTenant
var ContextKeyUserID = logging.ContextKeyUserID
var ContextKeyConsumer = logging.ContextKeyConsumer

func init() {
  configuration := logging.Configuration{
    Environment:  "development",
    Service:      "plutarch",
    Team:         "surveycorps",
    Project:      "te2eus",
    ConsoleLevel: "info",
    Version:      "1",
  }

  genericLog = logging.NewLogger(configuration)
  contextLog = logging.NewContextLogger(configuration)
}

func Info(
	message string,
	action string,
	payload map[string]string,
) {
	genericLog.Info(message, action, payload)
}

func Warn(
	message string,
	action string,
	payload map[string]string,
) {
	genericLog.Warn(message, action, payload)
}

func Debug(
	message string,
	action string,
	payload map[string]string,
) {
	genericLog.Debug(message, action, payload)
}

func Error(
	message string,
	action string,
	payload map[string]string,
) {
	genericLog.Error(message, action, payload)
}

func Fatal(
	message string,
	action string,
	payload map[string]string,
) {
	genericLog.Fatal(message, action, payload)
}

func CtxDebug(
	ctx context.Context,
	message string,
	action string,
	payload map[string]string,
) {
	contextLog.Debug(ctx, message, action, payload)
}

func CtxInfo(
	ctx context.Context,
	message string,
	action string,
	payload map[string]string,
) {
	contextLog.Info(ctx, message, action, payload)
}

func CtxWarn(
	ctx context.Context,
	message string,
	action string,
	payload map[string]string,
) {
	contextLog.Warn(ctx, message, action, payload)
}

func CtxError(
	ctx context.Context,
	message string,
	action string,
	payload map[string]string,
) {
	contextLog.Error(ctx, message, action, payload)
}

func CtxFatal(
	ctx context.Context,
	message string,
	action string,
	payload map[string]string,
) {
	contextLog.Fatal(ctx, message, action, payload)
}
```

You can use `os.Getenv("ENV_NAME")` to get the configuration value dinamically from the environment vars.

**NOTE: Do not forget to execute the command `go mod tidy` to install the library in the project.**

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

#### Version

Version of the service.

#### General

If any of the options is left empty it will not appear in the log and the default value for `ConsoleLevel` is `debug`.

## Usage examples

For all cases, the following fields will be automatically added to the log:

- method: from where the log was called.
- pid: proccess identifier.
- host: host name.

### Generic logger usage

This applies to all log types, the only difference is the log level.

Example with all the data:

```go
logger.Info("We are ready to GO!", "Starting server", map[string]string{
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
  "version": "1",
  "action": "Starting server",
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
  "version": "1"
}
```

The payload and action will be omitted if it is empty, the message will appear empty.

Additionally, the payload has a function that parses certain fields, currently they are:

- duration

Example:

If the duration field is sent with a string of numbers, it will be parsed to an int.

**NOTE: The duration field must be a string of only numbers, example: "9876", if it contains letters like "9876 ms" the string can't be parsed and the field will be returned with the value 0, this is because the duration field is used to measure and calculate response times and must be an integer.**

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
  "version": "1",
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
  "version": "1",
  "payload": { "duration": 0 }
}
```

### Context logger usage

It's the same as the generic log, the only difference is that it allows you to pass context as first argument to add extra values to the log.

Allowed context values:

- correlation_id
- causation_id
- tenant
- user_id
- consumer
- Datadog trace_id
- Datadog span_id

#### Correlation id

A unique id for each request that must be passed to every system that processes this request. Logging this id will make it easier to find related logs across different systems/files etc.

#### Causation id

You can also use an id that determine a correct ordering of the events that happend in you system. [Read more](https://blog.arkency.com/correlation-id-and-causation-id-in-evented-systems).

#### Tenant

Name of the client/sub-domain.

#### User id

It will facilitate investigating if user creates an incident ticket.

#### Consumer

Name of the consumer service.

#### Datadog trace id and span id

Those values are injected to the request context by instrumenting Datadog and is used to correlate traces with logs.

**NOTE: Any other context value will be ignored.**

The library provides the context type keys that you must use to match the values with the log.

Types:

- ContextKeyCorrelationID
- ContextKeyCausationID
- ContextKeyTenant
- ContextKeyUserID
- ContextKeyConsumer

**NOTE: Datadog's context keys are not provided because they are getting from their own span context, you don't have to do anything other than instrumenting your application with Datadog, the library will try to get those values for you.**

#### How to inject values to context

Because each framework operates in a different way, we don't not provided general functions or middlewares so it's up to you how to inject the values to the context.

Anyways here are examples that demostrate how you could do it.

Example with gin framework middleware:

```go
func loggerMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    var correlationId string
    var causationId string

    if id, err := uuid.Parse(c.GetHeader("correlation_id")); err == nil && c.GetHeader("correlation_id") != "" {
      correlationId = id.String()
      causationId = uuid.New().String()
    } else {
      correlationId = uuid.New().String()
      causationId = correlationId
    }

    ctx := c.Request.Context()

    ctx = context.WithValue(ctx, logger.ContextKeyCorrelationID, correlationId)
    ctx = context.WithValue(ctx, logger.ContextKeyCausationID, causationId)

    consumer := c.GetHeader("consumer")

    if consumer != "" {
      ctx = context.WithValue(ctx, logger.ContextKeyConsumer, consumer)
    }

    if value, exists := c.Get("session"); exists {
      session := value.(domain.Session)
      tenant, _ := session.Destructure()

      ctx = context.WithValue(ctx, logger.ContextKeyTenant, tenant)
      ctx = context.WithValue(ctx, logger.ContextKeyUserID, session.ID)
    }

    c.Request = c.Request.WithContext(ctx)

    c.Next()
  }
}
```

Router

```go
router.Use(loggerMiddleware())
...endpoints
```

Example with fiber framework middleware:

```go
func loggerMiddleware(c *fiber.Ctx) error {
  // TODO: implementing

  return c.Next()
}
```

Router

```go
router.Use(loggerMiddleware)
...endpoints
```

Assuming that all values were injected to the context, when using any log the values will be automatically added to the log.

Example:

```go
logger.CtxError(
  c.Request.Context(),
  "Parsing schema",
  err.Error(),
  map[string]string{"status": "500"},
)
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
  "version": "1",
  "correlation_id": "a6257245-ce50-4686-bfd1-668700fed150",
  "causation_id": "a6257245-ce50-4686-bfd1-668700fed150",
  "consumer": "pirithous",
  "dd": {
    "env": "development",
    "service": "plutarch",
    "span_id": "3890932795953018162",
    "trace_id": "3890932795953018162",
    "version": "1"
  },
  "ddsource": "go",
  "action": "Parsing schema",
  "payload": { "status": "500" }
}
```
