package plugins

import (
	"context"
	"errors"
	"fmt"
	"github.com/configwizard/gui/pkg/plugins/interop"
	"github.com/google/uuid"
	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const maxChunkSize = 1 * 1024 // 512 KB
const EXT = ".wasm"

type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	PluginId    string `json:"pluginId"`
}

type Plugin struct {
	Info Info
	interop.PluginService
}
type Manager struct {
	ctx     context.Context
	Plugins map[string]Plugin
}

func Init(pluginPath string, logger *log.Logger) *Manager {
	ctx := context.Background()

	// Initialize a Plugin loader
	p, err := interop.NewPluginServicePlugin(ctx)
	if err != nil {
		logger.Fatal("error starting Plugin service ", err)
	}

	var m Manager
	m.ctx = ctx
	m.Plugins = make(map[string]Plugin)
	filepath.Walk(pluginPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == EXT {
			if this, err := p.Load(ctx, path, HostFunctions{}); err != nil {
				return err
			} else {
				pluginTemporaryID := uuid.New().String()
				if pluginInfo, err := validatePlugin(this, pluginTemporaryID); err != nil {
					fmt.Println("could not validate Plugin ", err)
				} else {
					i := Info{
						Name:        pluginInfo.Name,
						Description: pluginInfo.Description,
						Author:      pluginInfo.Author,
						Version:     pluginInfo.Version,
						PluginId:    pluginInfo.PluginId,
					}

					m.Plugins[pluginTemporaryID] = Plugin{
						Info:          i,
						PluginService: this,
					}
				}
			}
		}
		return nil
	})
	//now get the Plugin content
	for _, v := range m.Plugins {
		content, err := v.Request(ctx, &interop.DataMessage{
			Type: interop.MessageType_CONTENT_REQUEST,
			Text: "content.html",
		})
		if err != nil {
			log.Fatal("error retrieving content ", err)
		}
		//now fire an event to see if the plugin received it
		fmt.Println("content ", content)
	}
	return &m
}

func validatePlugin(client interop.PluginService, pluginID string) (*interop.PluginInfo, error) {
	ctx := context.Background()
	info, err := client.InitializePlugin(ctx, &interop.PluginInfo{
		PluginId: pluginID,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("info ", info)
	// Perform your validation logic here
	if info.Name == "" {
		return nil, errors.New("Plugin name is required")
	}
	if info.Description == "" {
		return nil, errors.New("Plugin description is required")
	}
	if info.Author == "" {
		//check that the author is a registered author and Plugin is registered Plugin
		return nil, errors.New("Plugin author is required")
	}
	if info.Version == "" {
		return nil, errors.New("Plugin version is required")
	}
	return info, nil
}
func (a *Manager) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
	if p, ok := a.Plugins[pluginID]; !ok {
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

func (a *Manager) RequestPlugins() []Info {
	var p []Info
	fmt.Println("plugins requested ", a.Plugins)
	for _, v := range a.Plugins {
		p = append(p, v.Info)
	}
	return p
}

// myHostFunctions implements greeting.HostFunctions
type HostFunctions struct{}

// HttpGet is embedded into the Plugin and can be called by the Plugin.
func (HostFunctions) SignPayload(ctx context.Context, request *interop.DataMessage) (*interop.DataMessage, error) {
	return &interop.DataMessage{
		Type: 0,
		Text: "",
		Data: nil,
	}, nil
}

// Log is embedded into the Plugin and can be called by the Plugin.
func (HostFunctions) HostLog(ctx context.Context, request *interop.LogRequest) (*emptypb.Empty, error) {
	// Use the host logger
	log.Println("logging ", request.GetMessage())
	return &emptypb.Empty{}, nil
}

func (HostFunctions) Containers(ctx context.Context, _ *emptypb.Empty) (*interop.ElementsResponse, error) {
	var elements []*interop.Element
	// Example: Create 10 elements
	for i := 0; i < 10; i++ {
		elements = append(elements, &interop.Element{Id: "Container " + strconv.Itoa(i)})
	}
	return &interop.ElementsResponse{Elements: elements}, nil
}
func (HostFunctions) Container(ctx context.Context, el *interop.Element) (*interop.Element, error) {
	return &interop.Element{Id: el.Id}, nil
}
func (HostFunctions) Objects(ctx context.Context, _ *emptypb.Empty) (*interop.ElementsResponse, error) {
	var elements []*interop.Element
	// Example: Create 10 elements
	for i := 0; i < 10; i++ {
		elements = append(elements, &interop.Element{Id: "Object " + strconv.Itoa(i)})
	}
	return &interop.ElementsResponse{Elements: elements}, nil
}
func (HostFunctions) Object(ctx context.Context, el *interop.Element) (*interop.Element, error) {
	return &interop.Element{Id: el.Id}, nil
}
func (m Manager) RequestSign(ctx context.Context, request *interop.DataMessage) (interop.DataMessage, error) {
	return interop.DataMessage{
		Type: 0,
		Text: "",
		Data: nil,
	}, nil
}
