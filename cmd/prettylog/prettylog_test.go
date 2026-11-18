package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

func TestProcessInputLongLine(t *testing.T) {
	line := `{"level":"info","message":"` + strings.Repeat("a", 100*1024) + `"}`
	var out bytes.Buffer
	if err := processInput(strings.NewReader(line+"\n"+line+"\n"), zerolog.ConsoleWriter{Out: &out, NoColor: true}); err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(out.String(), "\n"); got != 2 {
		t.Fatalf("got %d lines, want 2", got)
	}
}
