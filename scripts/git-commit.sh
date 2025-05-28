#!/bin/bash

# Script to commit all Galaxy Node Pool changes following the Gelato Protocol standards

# Check for uncommitted changes
echo "Checking for changes..."
git status

# Add all files
echo "Adding all files..."
git add .

# Create commit message with proper format
echo "Creating commit message..."
COMMIT_MSG="(AI-ID: GALAXY-NODE-POOL-25052028) [feat]: Initial Galaxy Node Pool implementation"
COMMIT_DETAILS="
- Implemented modular, plugin-ready pool architecture
- Added .gal documentation following Gelato Protocol standards
- Created CLI for pool and node management
- Set up main net federation and incentives
- Organized documentation by component (pool, node) and purpose (versioning, phases, beta)
"

# Commit changes
echo "Committing changes..."
git commit -m "$COMMIT_MSG" -m "$COMMIT_DETAILS"

# Push changes (uncomment when ready)
# echo "Pushing changes..."
# git push origin main

echo "Git operations completed successfully!"
