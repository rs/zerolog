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

func main() {
	writer := zerolog.NewConsoleWriter()

	if isInputFromPipe() {
		scanner := bufio.NewScanner(os.Stdin)
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
	} else if len(os.Args) > 1 {
		for _, filename := range os.Args[1:] {
			// Scan each line from filename and write it into writer
			r, err := os.Open(filename)
			if err != nil {
				fmt.Printf("%s open: %v", filename, err)
				os.Exit(1)
			}
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				_, err = writer.Write(scanner.Bytes())
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					fmt.Printf("%s write: %v", filename, err)
					os.Exit(1)
				}
			}
			if err := scanner.Err(); err != nil {
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
