package plugins

import (
	"context"
	"errors"
	"fmt"
	"github.com/configwizard/gui/pkg/plugins/interop"
	"github.com/google/uuid"
	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

//func (m Manager) PluginEvent(ctx context.Context, request *interop.DataMessage) error {
//	for _, v := range m.Plugins {
//		v.PluginEvent(ctx, request)
//	}
//	return nil
//}

// myHostFunctions implements greeting.HostFunctions
type HostFunctions struct{}

func (m HostFunctions) PluginEvent(ctx context.Context, request *interop.DataMessage) (*emptypb.Empty, error) {
	//eventData := map[string]interface{}{
	//	"pluginID": request.Id,
	//	"data":     make([]byte, 0),
	//}
	runtime.EventsEmit(ctx, "plugin_backend_event", request)
	return &emptypb.Empty{}, nil
}

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
