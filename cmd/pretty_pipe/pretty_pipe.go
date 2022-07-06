package main

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func main() {
	if !isInputFromPipe() {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: app_with_zerolog |  2> >(pretty_pipe)")
		os.Exit(1)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	writer := zerolog.NewConsoleWriter()

	for scanner.Scan() {
		line := scanner.Text()

		_, err := writer.Write([]byte(line))

		if err != nil {
			println(line)
		}
	}
}
