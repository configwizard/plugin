package main

import (
	"context"
	"fmt"
	"github.com/configwizard/gui/pkg/plugins"
	"github.com/configwizard/gui/pkg/plugins/interop"
	"github.com/configwizard/sdk/controller"
	"github.com/configwizard/sdk/emitter"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
	return &Model{
		pluginManager: manager,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *Model) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *Model) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Trim leading slash and split URL
	pathParts := strings.Split(strings.TrimPrefix(req.URL.Path, "/"), "/")
	// Check if there's at least two parts (UUID and file path)
	if len(pathParts) < 2 {
		http.Error(res, "Invalid request format", http.StatusBadRequest)
		return
	}
	pluginID := pathParts[0]
	requestedFile := strings.Join(pathParts[1:], "/")

	// Check if the URL is correctly formatted
	if len(pathParts) < 2 {
		http.Error(res, "Invalid request format", http.StatusBadRequest)
		return
	}
	println("Requesting file:", requestedFile)
	//now we request the content
	if p, ok := a.pluginManager.Plugins[pluginID]; !ok {
		http.Error(res, "Invalid plugin", http.StatusBadRequest)
		return
	} else {
		contentResponse, err := p.Request(a.ctx, &interop.DataMessage{
			Type: interop.MessageType_CONTENT_REQUEST,
			Text: requestedFile,
		})
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFile)))
			return
		}
		// If the file is an HTML file, inject the base tag
		if strings.HasSuffix(requestedFile, ".html") {
			content := string(contentResponse.Data)
			// Inject the base tag after the opening head tag
			content = strings.Replace(content, "<head>", "<head>\n<base href=\"/"+pluginID+"/\">", 1)
			contentResponse.Data = []byte(content)
		}
		mimeType := "application/octet-stream" // Default MIME type if unknown
		if ext := filepath.Ext(requestedFile); ext != "" {
			if detectedType := mime.TypeByExtension(ext); detectedType != "" {
				mimeType = detectedType
			}
		}
		res.Header().Set("Content-Type", mimeType)
		res.Write(contentResponse.Data)
	}
}

func (a *Model) RequestPlugins() []plugins.Info {
	var p []plugins.Info
	fmt.Println("plugins requested ", a.pluginManager.Plugins)
	for _, v := range a.pluginManager.Plugins {
		p = append(p, v.Info)
	}
	return p
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
