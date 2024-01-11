
#https://github.com/knqyf263/go-plugin
generate:
	protoc -I pkg/plugins/interop/service --go-plugin_out=pkg/plugins/interop --go-plugin_opt=paths=source_relative pkg/plugins/interop/service/interop.proto

	#protoc -I pkg/plugins/interop/service --go-plugin_out=pkg/plugins/interop --go-plugin_opt=paths=source_relative pkg/plugins/interop/service/interop.proto

gen-example-plugin:
	tinygo build -o pkg/plugins/examples/basic.wasm -scheduler=none -target=wasi --no-debug pkg/plugins/examples/basic.go

