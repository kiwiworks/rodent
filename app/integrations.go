package app

import (
	"runtime"

	"github.com/kiwiworks/rodent/internal/golang"
	"github.com/kiwiworks/rodent/system/opt"
)

func moduleNameFromCallSite() string {
	return ExportedModuleNameFromCallSite()
}

// ExportedModuleNameFromCallSite is an exported version of moduleNameFromCallSite
// for testing purposes
func ExportedModuleNameFromCallSite() string {
	_, file, _, ok := runtime.Caller(2)
	if !ok {
		panic("failed to get caller while creating module, this is unexpected")
	}
	packageName, err := golang.FindModulePath(file)
	if err != nil {
		packageName = file
	}
	return packageName
}

func NewModule(opts ...opt.Option[Module]) Module {
	return NewNamedModule(
		moduleNameFromCallSite(),
		opts...,
	)
}
