package analyzer

import (
	"flag"
	"fmt"
)

type Config struct {
	IgnoredPkgs      ArrayFlag
	IgnoredFiles     ArrayFlag
	AllowedFinishers ArrayFlag
}

func DefaultConfig() Config {
	return Config{
		AllowedFinishers: []string{"Msg", "Msgf"},
	}
}

func FromFlags() Config {
	cfg := DefaultConfig()
	flag.Var(&cfg.IgnoredPkgs, "ignorePkg", "ignore the specified package")
	flag.Var(&cfg.IgnoredFiles, "ignoreFile", "ignore the specified file by its path and/or go path (package/file.go)")
	flag.Var(&cfg.AllowedFinishers, "finisher", "allowed finisher for the event chain")
	flag.Parse()

	return cfg
}

type ArrayFlag []string

func (i *ArrayFlag) String() string {
	return fmt.Sprintf("%v", []string(*i))
}

func (i *ArrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}
