package main

import (
	"os"

	"github.com/goplugin/pluginv3.0/core"
)

//go:generate make modgraph
func main() {
	os.Exit(core.Main())
}
