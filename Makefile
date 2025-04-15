.PHONY: serve clean release wasm-build

# Build WASM binary
wasm-build:
	tinygo build -o web/randomness.wasm -target wasm --no-debug web/main.go
	cp $(shell tinygo env TINYGOROOT)/targets/wasm_exec.js web/

# Serve the web directory
serve:
	go run contrib/server/server.go -webroot=web

# Clean build artifacts
clean:
	rm -f web/randomness.wasm web/wasm_exec.js

# Release a new version
release:
	./scripts/release.sh 