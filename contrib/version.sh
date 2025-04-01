#!/bin/bash

# Get the current version from version.go
current_version=$(grep "const Version" version.go | cut -d'"' -f2)

# Split version into major.minor.patch
IFS='.' read -r -a version_parts <<< "${current_version#v}"
major="${version_parts[0]}"
minor="${version_parts[1]}"
patch="${version_parts[2]}"

# Increment patch version
new_patch=$((patch + 1))
new_version="v${major}.${minor}.${new_patch}"

# Update version.go
sed -i '' "s/const Version = \"${current_version}\"/const Version = \"${new_version}\"/" version.go

# Print the new version
echo "${new_version}" 