// Package analyzer contains zerologlint implementation.
package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzer(config Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "zerologlint",
		Doc:  "This is a basic linter that checks for missing log event finishers.",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			t := getEvent(pass)
			if t == nil {
				return nil, nil // skip package, zerolog not used
			}
			cb(pass, config, t)
			return nil, nil
		},
	}
}

func cb(pass *analysis.Pass, config Config, eventType types.Type) {
	selections := getSelectionsWithReceiverType(pass, eventType)
	for _, selection := range selections {
		if hasBadFinishers(pass, config, selection) {
			pass.Reportf(selection.sel.Pos(), "missing or bad finisher for log chain, last call: %q", selection.fn.Name())
		}
	}
}

func hasBadFinishers(pass *analysis.Pass, config Config, s selection) bool {
	for _, ignoredPkg := range config.IgnoredPkgs {
		if s.pkg.Path() == ignoredPkg {
			return false
		}
	}

	for _, ignoredFile := range config.IgnoredFiles {
		if pass.Fset.File(s.fn.Pos()).Name() == ignoredFile {
			return false
		}
	}

	for _, f := range config.AllowedFinishers {
		if f == s.fn.Name() {
			return false
		}
	}

	return true
}
