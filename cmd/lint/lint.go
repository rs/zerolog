package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	recursivelyIgnoredPkgs arrayFlag
	ignoredPkgs            arrayFlag
	ignoredFiles           arrayFlag
	allowedFinishers       arrayFlag = []string{"Msg", "Msgf"}
	rootPkg                string
)

// parse input flags and args
func init() {
	flag.Var(&recursivelyIgnoredPkgs, "ignorePkgRecursively", "ignore the specified package and all subpackages recursively")
	flag.Var(&ignoredPkgs, "ignorePkg", "ignore the specified package")
	flag.Var(&ignoredFiles, "ignoreFile", "ignore the specified file by its path and/or go path (package/file.go)")
	flag.Var(&allowedFinishers, "finisher", "allowed finisher for the event chain")
	flag.Parse()

	// add zerolog to recursively ignored packages
	recursivelyIgnoredPkgs = append(recursivelyIgnoredPkgs, "github.com/rs/zerolog")
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "you must provide exactly one package path")
		os.Exit(1)
	}
	rootPkg = args[0]
}

var mode = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedDeps | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes | packages.NeedImports

func main() {
	// load the package and all its dependencies
	ps, err := packages.Load(&packages.Config{
		Mode: mode,
	}, rootPkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: unable to load the root package. %s\n", err.Error())
		os.Exit(1)
	}
	if len(ps) != 1 {
		fmt.Fprintf(os.Stderr, "Error: unable to load the root package: go %d root packages, expected 1 %s\n", len(ps), err.Error())
		os.Exit(1)
	}
	p := ps[0]
	// get the github.com/rs/zerolog.Event type
	event := getEvent(p)
	if event == nil {
		fmt.Fprintln(os.Stderr, "Error: github.com/rs/zerolog.Event declaration not found, maybe zerolog is not imported in the scanned package?")
		os.Exit(1)
	}

	// get all selections (function calls) with the github.com/rs/zerolog.Event (or pointer) receiver
	selections := getSelectionsWithReceiverType(p, event)

	// print the violations (if any)
	hasViolations := false
	for _, s := range selections {
		if hasBadFinisher(p, s) {
			hasViolations = true
			fmt.Printf("Error: missing or bad finisher for log chain, last call: %q at: %s:%v\n", s.fn.Name(), p.Fset.File(s.Pos()).Name(), p.Fset.Position(s.Pos()).Line)
		}
	}

	// if no violations detected, return normally
	if !hasViolations {
		fmt.Println("No violations found")
		return
	}

	// if violations were detected, return error code
	os.Exit(1)
}

func getEvent(p *packages.Package) types.Type {
	if strings.HasSuffix(p.PkgPath, "github.com/rs/zerolog") {
		for _, d := range p.TypesInfo.Defs {
			if d != nil && d.Name() == "Event" {
				return d.Type()
			}
		}
	}
	for _, imp := range p.Imports {
		event := getEvent(imp)
		if event != nil {
			return event
		}
	}
	return nil
}

func getSelectionsWithReceiverType(p *packages.Package, targetType types.Type) map[token.Pos]selection {
	selections := map[token.Pos]selection{}
	pre := func(current *packages.Package) bool {
		for i, t := range current.TypesInfo.Selections {
			switch o := t.Obj().(type) {
			case *types.Func:
				// this is not a bug, o.Type() is always *types.Signature, see docs
				if vt := o.Type().(*types.Signature).Recv(); vt != nil {
					typ := vt.Type()
					if pointer, ok := typ.(*types.Pointer); ok {
						typ = pointer.Elem()
					}

					if typ == targetType {
						if s, ok := selections[i.Pos()]; !ok || i.End() > s.End() {
							selections[i.Pos()] = selection{i, o, current.Types}
						}
					}
				}
			default:
				// skip
			}
		}
		return true
	}
	packages.Visit([]*packages.Package{p}, pre, nil)

	return selections
}

func hasBadFinisher(p *packages.Package, s selection) bool {
	pkgPath := strings.TrimPrefix(s.pkg.Path(), rootPkg+"/vendor/")
	absoluteFilePath := strings.TrimPrefix(p.Fset.File(s.Pos()).Name(), rootPkg+"/vendor/")
	goFilePath := pkgPath + "/" + filepath.Base(p.Fset.Position(s.Pos()).Filename)

	for _, f := range allowedFinishers {
		if f == s.fn.Name() {
			return false
		}
	}

	for _, ignoredPkg := range recursivelyIgnoredPkgs {
		if strings.HasPrefix(pkgPath, ignoredPkg) {
			return false
		}
	}

	for _, ignoredPkg := range ignoredPkgs {
		if pkgPath == ignoredPkg {
			return false
		}
	}

	for _, ignoredFile := range ignoredFiles {
		if absoluteFilePath == ignoredFile {
			return false
		}

		if goFilePath == ignoredFile {
			return false
		}
	}

	return true
}

type arrayFlag []string

func (i *arrayFlag) String() string {
	return fmt.Sprintf("%v", []string(*i))
}

func (i *arrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type selection struct {
	*ast.SelectorExpr
	fn  *types.Func
	pkg *types.Package
}
