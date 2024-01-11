//go:build tinygo.wasm

// Code generated by protoc-gen-go-plugin. DO NOT EDIT.
// versions:
// 	protoc-gen-go-plugin 0.8.0
// 	protoc               v4.25.1
// source: interop.proto

package interop

import (
	context "context"
	emptypb "github.com/knqyf263/go-plugin/types/known/emptypb"
	wasm "github.com/knqyf263/go-plugin/wasm"
	_ "unsafe"
)

const PluginServicePluginAPIVersion = 1

//export plugin_service_api_version
func _plugin_service_api_version() uint64 {
	return PluginServicePluginAPIVersion
}

var pluginService PluginService

func RegisterPluginService(p PluginService) {
	pluginService = p
}

//export plugin_service_initialize_plugin
func _plugin_service_initialize_plugin(ptr, size uint32) uint64 {
	b := wasm.PtrToByte(ptr, size)
	req := new(PluginInfo)
	if err := req.UnmarshalVT(b); err != nil {
		return 0
	}
	response, err := pluginService.InitializePlugin(context.Background(), req)
	if err != nil {
		ptr, size = wasm.ByteToPtr([]byte(err.Error()))
		return (uint64(ptr) << uint64(32)) | uint64(size) |
			// Indicate that this is the error string by setting the 32-th bit, assuming that
			// no data exceeds 31-bit size (2 GiB).
			(1 << 31)
	}

	b, err = response.MarshalVT()
	if err != nil {
		return 0
	}
	ptr, size = wasm.ByteToPtr(b)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

//export plugin_service_request
func _plugin_service_request(ptr, size uint32) uint64 {
	b := wasm.PtrToByte(ptr, size)
	req := new(DataMessage)
	if err := req.UnmarshalVT(b); err != nil {
		return 0
	}
	response, err := pluginService.Request(context.Background(), req)
	if err != nil {
		ptr, size = wasm.ByteToPtr([]byte(err.Error()))
		return (uint64(ptr) << uint64(32)) | uint64(size) |
			// Indicate that this is the error string by setting the 32-th bit, assuming that
			// no data exceeds 31-bit size (2 GiB).
			(1 << 31)
	}

	b, err = response.MarshalVT()
	if err != nil {
		return 0
	}
	ptr, size = wasm.ByteToPtr(b)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

//export plugin_service_process_data_stream
func _plugin_service_process_data_stream(ptr, size uint32) uint64 {
	b := wasm.PtrToByte(ptr, size)
	req := new(DataStreamMessage)
	if err := req.UnmarshalVT(b); err != nil {
		return 0
	}
	response, err := pluginService.ProcessDataStream(context.Background(), req)
	if err != nil {
		ptr, size = wasm.ByteToPtr([]byte(err.Error()))
		return (uint64(ptr) << uint64(32)) | uint64(size) |
			// Indicate that this is the error string by setting the 32-th bit, assuming that
			// no data exceeds 31-bit size (2 GiB).
			(1 << 31)
	}

	b, err = response.MarshalVT()
	if err != nil {
		return 0
	}
	ptr, size = wasm.ByteToPtr(b)
	return (uint64(ptr) << uint64(32)) | uint64(size)
}

type hostService struct{}

func NewHostService() HostService {
	return hostService{}
}

//go:wasm-module env
//export sign_payload
//go:linkname _sign_payload
func _sign_payload(ptr uint32, size uint32) uint64

func (h hostService) SignPayload(ctx context.Context, request *DataMessage) (*DataMessage, error) {
	buf, err := request.MarshalVT()
	if err != nil {
		return nil, err
	}
	ptr, size := wasm.ByteToPtr(buf)
	ptrSize := _sign_payload(ptr, size)
	wasm.FreePtr(ptr)

	ptr = uint32(ptrSize >> 32)
	size = uint32(ptrSize)
	buf = wasm.PtrToByte(ptr, size)

	response := new(DataMessage)
	if err = response.UnmarshalVT(buf); err != nil {
		return nil, err
	}
	return response, nil
}

//go:wasm-module env
//export host_log
//go:linkname _host_log
func _host_log(ptr uint32, size uint32) uint64

func (h hostService) HostLog(ctx context.Context, request *LogRequest) (*emptypb.Empty, error) {
	buf, err := request.MarshalVT()
	if err != nil {
		return nil, err
	}
	ptr, size := wasm.ByteToPtr(buf)
	ptrSize := _host_log(ptr, size)
	wasm.FreePtr(ptr)

	ptr = uint32(ptrSize >> 32)
	size = uint32(ptrSize)
	buf = wasm.PtrToByte(ptr, size)

	response := new(emptypb.Empty)
	if err = response.UnmarshalVT(buf); err != nil {
		return nil, err
	}
	return response, nil
}

//go:wasm-module env
//export containers
//go:linkname _containers
func _containers(ptr uint32, size uint32) uint64

func (h hostService) Containers(ctx context.Context, request *emptypb.Empty) (*ElementsResponse, error) {
	buf, err := request.MarshalVT()
	if err != nil {
		return nil, err
	}
	ptr, size := wasm.ByteToPtr(buf)
	ptrSize := _containers(ptr, size)
	wasm.FreePtr(ptr)

	ptr = uint32(ptrSize >> 32)
	size = uint32(ptrSize)
	buf = wasm.PtrToByte(ptr, size)

	response := new(ElementsResponse)
	if err = response.UnmarshalVT(buf); err != nil {
		return nil, err
	}
	return response, nil
}

//go:wasm-module env
//export container
//go:linkname _container
func _container(ptr uint32, size uint32) uint64

func (h hostService) Container(ctx context.Context, request *Element) (*Element, error) {
	buf, err := request.MarshalVT()
	if err != nil {
		return nil, err
	}
	ptr, size := wasm.ByteToPtr(buf)
	ptrSize := _container(ptr, size)
	wasm.FreePtr(ptr)

	ptr = uint32(ptrSize >> 32)
	size = uint32(ptrSize)
	buf = wasm.PtrToByte(ptr, size)

	response := new(Element)
	if err = response.UnmarshalVT(buf); err != nil {
		return nil, err
	}
	return response, nil
}

//go:wasm-module env
//export objects
//go:linkname _objects
func _objects(ptr uint32, size uint32) uint64

func (h hostService) Objects(ctx context.Context, request *emptypb.Empty) (*ElementsResponse, error) {
	buf, err := request.MarshalVT()
	if err != nil {
		return nil, err
	}
	ptr, size := wasm.ByteToPtr(buf)
	ptrSize := _objects(ptr, size)
	wasm.FreePtr(ptr)

	ptr = uint32(ptrSize >> 32)
	size = uint32(ptrSize)
	buf = wasm.PtrToByte(ptr, size)

	response := new(ElementsResponse)
	if err = response.UnmarshalVT(buf); err != nil {
		return nil, err
	}
	return response, nil
}

//go:wasm-module env
//export object
//go:linkname _object
func _object(ptr uint32, size uint32) uint64

func (h hostService) Object(ctx context.Context, request *Element) (*Element, error) {
	buf, err := request.MarshalVT()
	if err != nil {
		return nil, err
	}
	ptr, size := wasm.ByteToPtr(buf)
	ptrSize := _object(ptr, size)
	wasm.FreePtr(ptr)

	ptr = uint32(ptrSize >> 32)
	size = uint32(ptrSize)
	buf = wasm.PtrToByte(ptr, size)

	response := new(Element)
	if err = response.UnmarshalVT(buf); err != nil {
		return nil, err
	}
	return response, nil
}
