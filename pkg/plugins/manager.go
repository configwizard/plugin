package plugins

import (
	"context"
	"errors"
	"fmt"
	"github.com/configwizard/gui/pkg/plugins/interop"
	"github.com/google/uuid"
	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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
	m.Plugins = make(map[string]Plugin)
	filepath.Walk(pluginPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == EXT {
			if this, err := p.Load(ctx, path, HostFunctions{}); err != nil {
				return err
			} else {
				pluginTemporaryID := uuid.New().String()
				if pluginInfo, err := validatePlugin(this, pluginTemporaryID); err != nil {
					fmt.Println("could not validate Plugin")
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

//
//// HttpGet is embedded into the Plugin and can be called by the Plugin.
//func (m Manager) RequestSign(ctx context.Context, request *interop.DataMessage) (interop.DataMessage, error) {
//
//	fmt.Println("request to sign ", request.GetText())
//	return interop.DataMessage{
//		Type: 0,
//		Text: "",
//		Data: nil,
//	}, nil
//}
//
//// Log is embedded into the Plugin and can be called by the Plugin.
//func (m Manager) Log(ctx context.Context, request *interop.DataMessage) error {
//	// Use the host logger
//	log.Println("logging ", request.GetText())
//	return nil
//}

//
//func handleRetrieveNameData(cli interop.InteropServiceClient) string {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	r, err := cli.SendMessage(ctx, &interop.DataMessage{Type: interop.MessageType_NAME_REQUEST})
//	if err != nil {
//		log.Fatalf("could not greet: %v", err)
//	}
//	return r.GetText()
//}
//
//// handleStreamData handles the sending and receiving of data streams.
//func handleStreamData(cli interop.InteropServiceClient, buf *bytes.Buffer) {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	stream, err := cli.ProcessDataStream(ctx)
//	if err != nil {
//		log.Fatalf("Error creating data stream: %v", err)
//	}
//
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//
//	// Send data in a separate goroutine
//	go func() {
//		defer wg.Done()
//		randData := generateRandomData(1024 * 1024) // Generate 1MB of random data
//		for _, chunk := range chunkData(randData) {
//			if err := stream.Send(&interop.DataStreamMessage{Data: chunk}); err != nil {
//				log.Fatalf("Failed to send a chunk: %v", err)
//			}
//		}
//		stream.CloseSend()
//	}()
//
//	// Receive data
//	for {
//		in, err := stream.Recv()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			log.Fatalf("Failed to receive a chunk: %v", err)
//		}
//		processReceivedData(in.Data)
//	}
//	wg.Wait()
//}

// processReceivedData handles the data received from the server.
func processReceivedData(data []byte) {
	// Implement your data processing logic here
	fmt.Println("Received", len(data), "bytes of data")
}

// generateRandomData creates a slice of random bytes of a given length.
func generateRandomData(length int) []byte {
	data := make([]byte, length)
	rand.Read(data) // Note: This is not cryptographically secure
	return data
}

func chunkData(data []byte) [][]byte {
	var chunks [][]byte
	for len(data) > 0 {
		chunkSize := min(maxChunkSize, len(data))
		chunks = append(chunks, data[:chunkSize])
		data = data[chunkSize:]
	}
	return chunks
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
