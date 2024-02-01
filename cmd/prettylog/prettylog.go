package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func processInput(reader io.Reader) error {
	writer := zerolog.NewConsoleWriter()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		bytesToWrite := scanner.Bytes()
		_, err := writer.Write(bytesToWrite)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			fmt.Printf("%s\n", bytesToWrite)
		}
	}

	return scanner.Err()
}

func main() {

	if isInputFromPipe() {
		_ = processInput(os.Stdin)
	} else if len(os.Args) > 1 {
		for _, filename := range os.Args[1:] {
			// Scan each line from filename and write it into writer
			reader, err := os.Open(filename)
			if err != nil {
				fmt.Printf("%s open: %v", filename, err)
				os.Exit(1)
			}

			if err := processInput(reader); err != nil {
				fmt.Printf("%s scan: %v", filename, err)
				os.Exit(1)
			}
		}
	} else {
		fmt.Println("Usage:")
		fmt.Println("  app_with_zerolog | 2> >(prettylog)")
		fmt.Println("  prettylog zerolog_output.jsonl")
		os.Exit(1)
		return
	}
}
