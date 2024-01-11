package main

import (
	"github.com/configwizard/gui/pkg/plugins"
	"log"
	"path/filepath"
)

func main() {
	pluginPath := filepath.Join("/Users/alexwalker/go/src/github.com/configwizard/gui/pkg/plugins", "examples")
	plugins.Init(pluginPath, &log.Logger{})
}
