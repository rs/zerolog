package main

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func main() {
	if !isInputFromPipe() {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: app_with_zerolog |  2> >(prettylog)")
		os.Exit(1)
		return
	}

	writer := zerolog.NewConsoleWriter()
	_, _ = io.Copy(writer, os.Stdin)
}
