.PHONY: build
build:
	env GOOS=js GOARCH=wasm go build -o build/galaxy.wasm .
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" build/
	cp main.html build/index.html
