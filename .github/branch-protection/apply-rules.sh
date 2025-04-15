#!/bin/bash

set -e

# Function to handle errors
handle_error() {
    echo "Error: $1"
    exit 1
}

# Function to check if a command exists
check_command() {
    local cmd=$1
    if ! command -v "$cmd" >/dev/null 2>&1; then
        handle_error "Required command '$cmd' not found. Please install it before running this script."
    fi
}

# Function to compare versions
version_gt() {
    test "$(printf '%s\n' "$@" | sort -V | head -n 1)" != "$1"
}

# Check for required commands
check_command "gh"
check_command "yq"
check_command "jq"

# Check gh version (need 2.39.0 or higher for rulesets)
GH_VERSION=$(gh --version | head -n 1 | cut -d' ' -f3)
MIN_VERSION="2.39.0"

if version_gt "$MIN_VERSION" "$GH_VERSION"; then
    handle_error "gh version $GH_VERSION is too old. Please upgrade to version $MIN_VERSION or higher."
fi

# Get repository name
REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner) || handle_error "Failed to get repo name"
echo "Repository: $REPO"

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Function to get ruleset list
get_ruleset_list() {
    gh ruleset list || handle_error "Failed to list rulesets"
}

# Function to get ruleset ID by name
get_ruleset_id() {
    local name=$1
    local ruleset_list
    ruleset_list=$(get_ruleset_list)
    echo "$ruleset_list" | awk -v name="$name" '$2 == name {print $1}'
}

# Function to process template file
process_template() {
    local template_file=$1
    # Create a temporary file
    local temp_file=$(mktemp)
    # Use yq to replace the template variable and convert to JSON
    yq -o json "del(.precedence) | .source = \"$REPO\"" "$template_file" > "$temp_file"
    echo "$temp_file"
}

# Function to apply a branch protection rule
apply_rule() {
    local rule_file=$1
    echo "Applying rule from $rule_file..."
    
    # Process the template
    local processed_file
    processed_file=$(process_template "$rule_file")
    
    # Extract rule name for checking existing rulesets
    local rule_name=$(jq -r '.name' "$processed_file") || handle_error "Failed to extract rule name"
    
    # Check if ruleset already exists
    local ruleset_id
    ruleset_id=$(get_ruleset_id "$rule_name")
    
    if [ -n "$ruleset_id" ]; then
        echo "Ruleset '$rule_name' already exists (ID: $ruleset_id), updating..."
        # Update the existing ruleset
        gh api --method PUT \
            -H "Accept: application/vnd.github+json" \
            -H "X-GitHub-Api-Version: 2022-11-28" \
            "/repos/$REPO/rulesets/$ruleset_id" \
            --input "$processed_file" || handle_error "Failed to update ruleset"
    else
        echo "Creating new ruleset '$rule_name'..."
        # Create new ruleset
        gh api --method POST \
            -H "Accept: application/vnd.github+json" \
            -H "X-GitHub-Api-Version: 2022-11-28" \
            "/repos/$REPO/rulesets" \
            --input "$processed_file" || handle_error "Failed to create ruleset"
    fi
    
    # Clean up temporary file
    rm "$processed_file"
}

# Function to get precedence from YAML file
get_precedence() {
    local file=$1
    yq -r '.precedence // 999' "$file"
}

# Function to reorder rulesets
reorder_rulesets() {
    local ruleset_list
    ruleset_list=$(get_ruleset_list)
    echo "Current rulesets:"
    echo "$ruleset_list"
    
    # Delete all existing rulesets
    while IFS=$'\t' read -r id name _ _ _; do
        if [[ $id =~ ^[0-9]+$ ]]; then
            echo "Deleting ruleset $id ($name)..."
            gh api --method DELETE \
                -H "Accept: application/vnd.github+json" \
                -H "X-GitHub-Api-Version: 2022-11-28" \
                "/repos/$REPO/rulesets/$id" || handle_error "Failed to delete ruleset"
        fi
    done <<< "$ruleset_list"
    
    # Get all YAML files and sort by precedence
    local -a yaml_files
    for file in "$SCRIPT_DIR"/*.yaml; do
        if [ -f "$file" ]; then
            yaml_files+=("$file")
        fi
    done
    
    # Sort files by precedence
    IFS=$'\n' sorted_files=($(for file in "${yaml_files[@]}"; do
        echo "$(get_precedence "$file") $file"
    done | sort -n | cut -d' ' -f2-))
    
    # Apply rules in precedence order
    for file in "${sorted_files[@]}"; do
        apply_rule "$file"
    done
}

# Process all rules and ensure correct ordering
echo "Processing all ruleset files..."
reorder_rulesets

# Final verification
echo "Verifying all rulesets..."
get_ruleset_list 