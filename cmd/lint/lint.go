package main

import (
	"github.com/rs/zerolog/cmd/lint/analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzer(analyzer.FromFlags()))
}
