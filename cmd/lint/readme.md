# Zerolog Lint

This is a basic linter that checks for missing log event finishers. Finds errors like: `log.Error().Int64("userID": 5)` - missing the `Msg`/`Msgf` finishers.

## Problem

When using zerolog it's easy to forget to finish the log event chain by calling a finisher - the `Msg` or `Msgf` function that will schedule the event for writing. The problem with this is that it doesn't warn/panic during compilation and it's not easily found by grep or other general tools. It's even prominently mentioned in the project's readme, that:

> It is very important to note that when using the **zerolog** chaining API, as shown above (`log.Info().Msg("hello world"`), the chain must have either the `Msg` or `Msgf` method call. If you forget to add either of these, the log will not occur and there is no compile time error to alert you of this.

## Solution

A basic linter like this one here that looks for method invocations on `zerolog.Event` can examine the last call in a method call chain and check if it is a finisher, thus pointing out these errors.

## Usage

Just compile this and then run it. Or just run it via `go run` command via something like `go run cmd/lint/lint.go`.

The command accepts only one argument - the package to be inspected - and 4 optional flags, all of which can occur multiple times. The standard synopsis of the command is:

`lint [-finisher value] [-ignoreFile value] [-ignorePkg value] [-ignorePkgRecursively value] package`

#### Flags

- finisher
    - specify which finishers to accept, defaults to `Msg` and `Msgf`
- ignoreFile
    - which files to ignore, either by full path or by go path (package/file.go)
- ignorePkg
    - do not inspect the specified package if found in the dependecy tree
- ignorePkgRecursively
    - do not inspect the specified package or its subpackages if found in the dependency tree

## Drawbacks

As it is, linter can generate a false positives in a specific case. These false positives come from the fact that if you have a method that returns a `zerolog.Event` the linter will flag it because you are obviously not finishing the event. This will be solved in later release.

