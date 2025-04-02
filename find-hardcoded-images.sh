#!/bin/bash
# Search for hardcoded image data in the codebase

cd /Users/subrahmanyagonella/the-repo/sentinelstacks

echo "Searching for hardcoded image data..."
grep -r "sha256:abcde" --include="*.go" .

echo "Searching for user/chatbot..."
grep -r "user/chatbot" --include="*.go" .

echo "Searching for getSampleImages..."
grep -r "getSampleImages" --include="*.go" .

echo "Search complete."
