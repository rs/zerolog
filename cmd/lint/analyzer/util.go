package analyzer

import (
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func getEvent(pass *analysis.Pass) types.Type {
	for _, pkg := range pass.Pkg.Imports() {
		if strings.HasSuffix(pkg.Path(), "github.com/rs/zerolog") {
			d := pkg.Scope().Lookup("Event")
			if d != nil {
				return d.Type()
			}
		}
	}

	return nil
}
