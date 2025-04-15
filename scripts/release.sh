#!/bin/bash

set -e

# Create a release PR using release-please
npx release-please release-pr --package-name=randomness --token=$GITHUB_TOKEN 