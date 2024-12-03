define HTML_CONTENT
<!DOCTYPE html>
<script src="wasm_exec.js"></script>
<script>
const go = new Go();
WebAssembly.instantiateStreaming(fetch("yourgame.wasm"), go.importObject).then(result => {
    go.run(result.instance);
});
</script>
endef

export HTML_CONTENT
.PHONY: build
build:
	env GOOS=js GOARCH=wasm go build -o build/galaxy.wasm .
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" build/
	echo "$$HTML_CONTENT" > build/index.html
