#!/bin/bash

set -e

# Create release directory
RELEASE_DIR="release"
mkdir -p $RELEASE_DIR

# Get TinyGo version
TINYGO_VERSION=$(tinygo version | cut -d' ' -f2)

# Build WASM binary
echo "Building WASM binary..."
tinygo build -o $RELEASE_DIR/randomness.wasm \
    -target wasm \
    --no-debug \
    -opt=2 \
    -gc=leaking \
    -panic=trap \
    -scheduler=none \
    web/main.go

# Copy wasm_exec.js
echo "Copying wasm_exec.js..."
cp $(tinygo env TINYGOROOT)/targets/wasm_exec.js $RELEASE_DIR/

# Generate checksums
echo "Generating checksums..."
cd $RELEASE_DIR
sha256sum *.wasm *.js > checksums.txt

# Generate signatures (requires GPG key)
echo "Generating signatures..."
for file in *.wasm *.js; do
    gpg --armor --detach-sign "$file"
done

# Create release notes
echo "Creating release notes..."
cat > release_notes.md << EOF
# Release ${GITHUB_REF#refs/tags/}

## Build Information
- TinyGo Version: ${TINYGO_VERSION}

## Files
- randomness.wasm
- wasm_exec.js

## Checksums
\`\`\`
$(cat checksums.txt)
\`\`\`

## Signatures
- randomness.wasm.asc
- wasm_exec.js.asc
EOF

# Output release info for GitHub Actions
echo "release_dir=$RELEASE_DIR" >> $GITHUB_OUTPUT 