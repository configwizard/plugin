package main

import (
	"context"
	"embed"
	"errors"
	"github.com/configwizard/gui/pkg/plugins/interop"
	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"io/fs"
	"os"
	"path/filepath"
)

var assets fs.FS

//go:embed assets
var embeddedAssets embed.FS

var isDevMode string // This will be set based on the build argument

func retrievePageContent(fsys fs.FS, page string) ([]byte, error) {
	// Join the directory and file name, then read the file
	var path string = page
	if isDevMode != "true" {
		path = filepath.Join("assets", page)
	}
	content, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func init() {
	if isDevMode == "true" {
		// Development mode: load assets from disk
		assets = os.DirFS("path/to/assets") //change this to the absolute path of your assets.
	} else {
		// Production mode: use embedded assets
		assets = embeddedAssets
	}
}

// main is required for TinyGo to compile to Wasm.
func main() {
	interop.RegisterPluginService(&myPlugin{
		Name:        "my plugin",
		Description: "an example plugin",
		Author:      "alex walker",
		Version:     "v0.0.1",
	})
}

type myPlugin struct {
	host        interop.HostService
	ctx         context.Context
	Name        string
	Description string
	Author      string
	Version     string
	PluginID    string
}

func (m *myPlugin) InitializePlugin(ctx context.Context, info *interop.PluginInfo) (*interop.PluginInfo, error) {
	m.ctx = ctx
	m.PluginID = info.PluginId
	m.host = interop.NewHostService()
	return &interop.PluginInfo{
		Name:        m.Name,
		Description: m.Description,
		Author:      m.Author,
		Version:     m.Version,
		PluginId:    m.PluginID,
	}, nil
}

func (m myPlugin) PluginEvent(ctx context.Context, request *interop.DataMessage) (*emptypb.Empty, error) {

	m.host.HostLog(ctx, &interop.LogRequest{
		Message: "request data - " + string(request.Data),
	})
	switch request.GetText() {
	case "retrieveContainers":
		m.host.HostLog(ctx, &interop.LogRequest{
			Message: "internal_retrieve_Containers event - " + request.GetText(),
		})
		if _, err := m.host.PluginEvent(ctx, request); err != nil {
			m.host.HostLog(ctx, &interop.LogRequest{
				Message: "error sending - " + err.Error(),
			})
		}
	}
	return nil, nil
}
func (m myPlugin) Request(ctx context.Context, request *interop.DataMessage) (*interop.DataMessage, error) {
	m.host.HostLog(ctx, &interop.LogRequest{
		Message: "request - " + request.GetText(),
	})
	switch request.GetType() {
	//case interop.MessageType_INITIALIZE:
	//	hostFunctions.HostLog(ctx, &interop.LogRequest{
	//		Message: "initialise plugin",
	//	})
	//	return &interop.DataMessage{
	//		Type: interop.MessageType_INITIALIZE,
	//		Text: plugins.READY,
	//	}, nil
	case interop.MessageType_ICON_REQUEST:
		icon, err := retrievePageContent(assets, "icon.svg")
		if err != nil {
			return nil, err
		}
		return &interop.DataMessage{
			Type: interop.MessageType_ICON_REQUEST,
			Text: "icon.svg",
			Data: icon,
		}, nil
	case interop.MessageType_NAME_REQUEST:
		return &interop.DataMessage{
			Type: interop.MessageType_NAME_REQUEST,
			Text: m.Name,
			Data: nil,
		}, nil
	case interop.MessageType_DESCRIPTION_REQUEST:
		return &interop.DataMessage{
			Type: interop.MessageType_NAME_REQUEST,
			Text: m.Description,
			Data: nil,
		}, nil
	case interop.MessageType_CONTENT_REQUEST:
		page := request.GetText() //the page name to retrieve
		if page == "" {
			m.host.HostLog(ctx, &interop.LogRequest{
				Message: "need to set a page",
			})
			return nil, errors.New("need to supply a content path")
		}
		content, err := retrievePageContent(assets, page)
		if err != nil {
			return nil, err
		}
		return &interop.DataMessage{
			Type: interop.MessageType_CONTENT_REQUEST,
			Text: page,
			Data: content,
		}, nil
	case interop.MessageType_DATA_PROCESS_REQUEST:
	case interop.MessageType_UNKNOWN_REQUEST:
	case interop.MessageType_DEFAULT:
	default:
		return &interop.DataMessage{
			Type: interop.MessageType_UNKNOWN_REQUEST,
			Text: "unimplemented",
		}, nil
	}
	return request, errors.New("case not implemented")
}
func (m myPlugin) ProcessDataStream(ctx context.Context, str *interop.DataStreamMessage) (*interop.DataStreamMessage, error) {
	return &interop.DataStreamMessage{
		Data: make([]byte, 0),
	}, nil
}
