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

func main() {
	args := flag.Args()

	// load the package and all its dependencies
	conf := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: true,
	}
	pkgs, err := packages.Load(conf, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot load packages: %v\n", err)
		os.Exit(1)
	}

	// TODO do this with packages instead of loader -- there's a mapping to most of they types we need somewhere
	// get the github.com/rs/zerolog.Event type
	event := getEvent(pkgs)
	if event == nil {
		fmt.Fprintln(os.Stderr, "Error: github.com/rs/zerolog.Event declaration not found, maybe zerolog is not imported in the scanned package?")
		os.Exit(1)
	}

	// print the violations (if any)
	hasViolations := false

	// ensure one checks per package (dups may show up otherwise)
	var uniquePackages = make(map[string]*packages.Package)
	for _, pkg := range pkgs {
		uniquePackages[pkg.PkgPath] = pkg
	}

	for _, pkg := range uniquePackages {
		selections := getSelectionsWithReceiverType(pkg, event)

		for _, s := range selections {
			if hasBadFinisher(pkg, s) {
				hasViolations = true
				fmt.Printf("Error: missing or bad finisher for log chain, last call: %q at: %s:%v\n", s.fn.Name(), pkg.Fset.File(s.Pos()).Name(), pkg.Fset.Position(s.Pos()).Line)
			}
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

func getEvent(pkgs []*packages.Package) types.Type {
	var packagesToCheck = make(map[string]*packages.Package)
	for _, pkg := range pkgs {
		for _, imp := range pkg.Imports {
			packagesToCheck[imp.PkgPath] = imp
		}
		packagesToCheck[pkg.PkgPath] = pkg
	}

	for _, pkg := range packagesToCheck {
		if strings.HasSuffix(pkg.PkgPath, "github.com/rs/zerolog") {
			for _, d := range pkg.TypesInfo.Defs {
				if d != nil && d.Name() == "Event" {
					return d.Type()
				}
			}
		}
	}

	return nil
}

func getSelectionsWithReceiverType(pkg *packages.Package, targetType types.Type) map[token.Pos]selection {
	selections := map[token.Pos]selection{}
	for i, t := range pkg.TypesInfo.Selections {
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
						selections[i.Pos()] = selection{i, o, pkg.Types}
					}
				}
			}
		default:
			// skip
		}
	}

	return selections
}

func hasBadFinisher(pkg *packages.Package, s selection) bool {
	pkgPath := strings.TrimPrefix(s.pkg.Path(), rootPkg+"/vendor/")
	absoluteFilePath := strings.TrimPrefix(pkg.Fset.File(s.Pos()).Name(), rootPkg+"/vendor/")
	goFilePath := pkgPath + "/" + filepath.Base(pkg.Fset.Position(s.Pos()).Filename)

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
