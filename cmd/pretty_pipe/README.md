# Zerolog Pretty Pipe

This is a basic CLI utility that will colorize and pretty print your JSON logs.

## Problem

When using zerolog it's sometimes difficult to do an environment based log formatting. 
For production purposes it is useful to have JSON logs to pipe them through Kibana for efficient log entry analysis.

However, when using such applications in development environment there is little use for structured JSON logs.
Humans are visual beings so the preferred format should be ordinary `Console` format with colors if possible.

Zerolog ships with `ConsoleWriter` which does exactly that. But modifying your code to get pretty printed and colorized
logs is somewhat difficult.

My solution to do this programmatically was something like this:

```go
func init() {
	if !IsProductionEnv() {
		log.Logger = log.Output(zerolog.NewConsoleWriter())
	}
}

func GetLogger(component string) zerolog.Logger {
	return log.With().
		Str("Component", module).
		Logger()
}
```

Here `init` is called automatically before `main` by go. `GetLogger` function is used to initialize new loggers.
If you do something like this:
```go
func init() {
	if !IsProductionEnv() {
		log.Logger = log.Output(zerolog.NewConsoleWriter())
	}
}

var logger = log.With().Str("Component", "Module").Logger()
```

It won't work since global variables are initialized before `init` is called.

## Solution

A basic CLI tool to pipe application output through it to pretty print & colorize JSON logs to console friendly logs.

## Usage

You can compile it or run it directly. The only issue is that by default Zerolog does not output to `stdout` 
but rather to `stderr` so we must pipe `stderr` stream to this CLI tool.

### Linux

These commands will redirect `stderr` to our `pretty_pipe` tool and `stdout` will remain unaffected.

1. Compiled version
```shell
some_program_with_zerolog 2> >(pretty_pipe)
```

2. Run it directly with `go run`
```shell
some_program_with_zerolog 2> >(go run cmd/pretty_pipe/pretty_pipe.go)
```

### Windows

These commands will redirect `stderr` to `stdout` and then pipe it to our `pretty_pipe` tool.

1. Compiled version
```shell
some_program_with_zerolog 2>&1 | pretty_pipe 
```

2. Run it directly with `go run`
```shell
some_program_with_zerolog 2>&1 | go run cmd/pretty_pipe/pretty_pipe.go
```
