package utils

import (
	"runtime/debug"
)

const ModName = "github.com/selectel/private-dns-go"

var Version = version()

func version() string {
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range info.Deps {
			if dep.Path == ModName {
				return dep.Version
			}
		}
	}

	return "v0.0.0"
}
