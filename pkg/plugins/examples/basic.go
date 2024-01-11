package main

import (
	"context"
	"embed"
	"errors"
	"github.com/configwizard/gui/pkg/plugins/interop"
	"path/filepath"
)

//go:embed assets
var assets embed.FS

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
	Name        string
	Description string
	Author      string
	Version     string
	PluginID    string
}

func (m *myPlugin) InitializePlugin(ctx context.Context, info *interop.PluginInfo) (*interop.PluginInfo, error) {
	m.PluginID = info.PluginId
	return &interop.PluginInfo{
		Name:        m.Name,
		Description: m.Description,
		Author:      m.Author,
		Version:     m.Version,
		PluginId:    m.PluginID,
	}, nil
}
func (m myPlugin) Request(ctx context.Context, request *interop.DataMessage) (*interop.DataMessage, error) {

	hostFunctions := interop.NewHostService()
	hostFunctions.HostLog(ctx, &interop.LogRequest{
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
		icon, err := assets.ReadFile("assets/icon.svg")
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
			Text: "basic plugin",
			Data: nil,
		}, nil
	case interop.MessageType_DESCRIPTION_REQUEST:
		return &interop.DataMessage{
			Type: interop.MessageType_NAME_REQUEST,
			Text: "an example plugin",
			Data: nil,
		}, nil
	case interop.MessageType_CONTENT_REQUEST:
		page := request.GetText() //the page name to retrieve
		if page == "" {
			hostFunctions.HostLog(ctx, &interop.LogRequest{
				Message: "need to set a page",
			})
			return nil, errors.New("need to supply a content path")
		}
		content, err := assets.ReadFile(filepath.Join("assets/", page))
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
