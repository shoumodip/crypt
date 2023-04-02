.PHONY: all
all: crypt server web/crypt.wasm

crypt: cmd/crypt/crypt.go pkg/crypt/crypt.go pkg/field/field.go
	go build cmd/crypt/crypt.go

server: cmd/server/server.go
	go build cmd/server/server.go

web/crypt.wasm: cmd/wasm/wasm.go pkg/crypt/crypt.go pkg/field/field.go
	GOOS=js GOARCH=wasm go build -o web/crypt.wasm cmd/wasm/wasm.go
