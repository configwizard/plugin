package plugins

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestManager(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("error retrieving working directory ", err)
	}
	pluginPath := filepath.Join(wd, "examples")
	Init(pluginPath, &log.Logger{})
}
