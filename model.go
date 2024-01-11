package main

import (
	"context"
	"github.com/configwizard/gui/pkg/plugins"
	"github.com/configwizard/sdk/controller"
	"github.com/configwizard/sdk/emitter"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var version string

// App struct
type Model struct {
	http.Handler
	pluginManager *plugins.Manager
	ctx           context.Context
	controller    controller.Controller
}

// NewApp creates a new App application struct
func NewModel() *Model {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("error retrieving working directory ", err)
	}
	pluginPath := filepath.Join(wd, "pkg", "plugins", "examples")
	manager := plugins.Init(pluginPath, &log.Logger{})
	version = "basic-dev-env"
	return &Model{
		pluginManager: manager,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *Model) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *Model) RequestContainers() error {
	var c = struct {
		Name string
		Size int64
	}{
		Size: 10,
	}

	go func() {
		for i := 0; i < 10; i++ {
			c.Name = generateRandomString(8)
			runtime.EventsEmit(a.ctx, emitter.ContainerAddUpdate, c)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	return nil
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	lettersAndDigits := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRXYZ"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(lettersAndDigits[rand.Intn(len(lettersAndDigits))])
	}
	return sb.String()
}
