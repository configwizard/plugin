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

type Manager struct {
	Plugins map[string]interop.PluginService
}

func Init(pluginPath string, logger *log.Logger) *Manager {
	ctx := context.Background()

	// Initialize a plugin loader
	p, err := interop.NewPluginServicePlugin(ctx)
	if err != nil {
		logger.Fatal("error starting plugin service ", err)
	}

	var m Manager
	m.Plugins = make(map[string]interop.PluginService)
	filepath.Walk(pluginPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == EXT {
			if this, err := p.Load(ctx, path, HostFunctions{}); err != nil {
				return err
			} else {
				pluginTemporaryID := uuid.New().String()
				if err := validatePlugin(this, pluginTemporaryID); err != nil {
					fmt.Println("could not validate plugin")
				} else {
					m.Plugins[pluginTemporaryID] = this
				}
			}
		}
		return nil
	})
	//now get the plugin content
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

func validatePlugin(client interop.PluginService, pluginID string) error {
	ctx := context.Background()
	info, err := client.InitializePlugin(ctx, &interop.PluginInfo{
		PluginId: pluginID,
	})
	if err != nil {
		return err
	}
	fmt.Println("info ", info)
	// Perform your validation logic here
	if info.Name == "" {
		return errors.New("plugin name is required")
	}
	if info.Description == "" {
		return errors.New("plugin description is required")
	}
	if info.Author == "" {
		//check that the author is a registered author and plugin is registered plugin
		return errors.New("plugin author is required")
	}
	if info.Version == "" {
		return errors.New("plugin version is required")
	}

	return nil
}

// myHostFunctions implements greeting.HostFunctions
type HostFunctions struct{}

// HttpGet is embedded into the plugin and can be called by the plugin.
func (HostFunctions) SignPayload(ctx context.Context, request *interop.DataMessage) (*interop.DataMessage, error) {
	return &interop.DataMessage{
		Type: 0,
		Text: "",
		Data: nil,
	}, nil
}

// Log is embedded into the plugin and can be called by the plugin.
func (HostFunctions) HostLog(ctx context.Context, request *interop.LogRequest) (*emptypb.Empty, error) {
	// Use the host logger
	log.Println("logging ", request.GetMessage())
	return &emptypb.Empty{}, nil
}

//
//// HttpGet is embedded into the plugin and can be called by the plugin.
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
//// Log is embedded into the plugin and can be called by the plugin.
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
