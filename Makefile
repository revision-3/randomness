.PHONY: serve clean release wasm-build

# Build WASM binary
wasm-build:
	tinygo build -o web/wasm.wasm -target wasm --no-debug web/main.go
	cp $(shell tinygo env TINYGOROOT)/targets/wasm_exec.js web/

# Serve the web directory
serve:
	go run contrib/server/server.go -webroot=web

# Clean build artifacts
clean:
	rm -f web/wasm.wasm web/wasm_exec.js

# Release a new version
release:
	@new_version=$$(./contrib/version.sh) && \
	git add version.go && \
	git commit -m "Release $$new_version" && \
	git tag $$new_version && \
	git push origin main && \
	git push origin $$new_version 