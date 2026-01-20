package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var timeFormats map[string]string = map[string]string{
	"default":     time.Kitchen,
	"ansic":       time.ANSIC,
	"unix":        time.UnixDate,
	"rfc822":      time.RFC822,
	"rfc822z":     time.RFC822Z,
	"rfc850":      time.RFC850,
	"rfc1123":     time.RFC1123,
	"rfc1123z":    time.RFC1123Z,
	"rfc3339":     time.RFC3339,
	"rfc3339nano": time.RFC3339Nano,
	"stamp":       time.Stamp,
	"stampmilli":  time.StampMilli,
	"datetime":    time.DateTime,
	"timeonly":    time.TimeOnly,
	"full":        time.RFC1123,
}

func allowedTimeFormats() string {
	formats := make([]string, 0, len(timeFormats))
	for k := range timeFormats {
		formats = append(formats, k)
	}
	sort.Strings(formats)

	return strings.Join(formats, ", ")
}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func processInput(reader io.Reader, writer io.Writer) error {
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

func getTimeFormat(flagValue string) string {
	format, ok := timeFormats[flagValue]
	if !ok {
		return flagValue
	}

	return format
}

func main() {
	timeFormatFlag := flag.String(
		"time-format",
		"default",
		"Time format, one of: "+allowedTimeFormats()+" or a custom golang time format",
	)
	timeLocationFlag := flag.String(
		"time-location",
		"UTC",
		"Time location, one of: UTC, Local or a custom location",
	)

	flag.Parse()
	loc, err := time.LoadLocation(*timeLocationFlag)
	if err != nil {
		fmt.Printf("time location %s: %v", *timeLocationFlag, err)
		os.Exit(1)
	}

	writer := zerolog.NewConsoleWriter()
	writer.TimeFormat = getTimeFormat(*timeFormatFlag)
	writer.TimeLocation = loc

	if flag.NArg() >= 1 {
		for _, filename := range flag.Args() {
			// Scan each line from filename and write it into writer
			reader, err := os.Open(filename)
			if err != nil {
				fmt.Printf("%s open: %v", filename, err)
				os.Exit(1)
			}

			if err := processInput(reader, writer); err != nil {
				fmt.Printf("%s scan: %v", filename, err)
				os.Exit(1)
			}
		}
	} else if isInputFromPipe() {
		_ = processInput(os.Stdin, writer)
	} else {
		fmt.Println("Usage:")
		fmt.Println("  app_with_zerolog | 2> >(prettylog)")
		fmt.Println("  prettylog zerolog_output.jsonl")
		os.Exit(1)
		return
	}
}
