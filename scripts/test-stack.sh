#!/bin/bash
# Test script for SentinelStacks Stack functionality

set -e  # Exit on error

# Colors for output
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
RESET="\033[0m"

# Paths
SENTINEL_PATH="./sentinel"
EXAMPLES_DIR="./examples/stacks"
TEST_DIR="./test-stack-output"

echo -e "${BLUE}=== SentinelStacks Stack Engine Test ===${RESET}"

# Check if sentinel binary exists
if [ ! -f "$SENTINEL_PATH" ]; then
    echo -e "${RED}Error: sentinel binary not found at $SENTINEL_PATH${RESET}"
    echo -e "${YELLOW}Please build the binary first with 'make build'${RESET}"
    exit 1
fi

# Create test directory
mkdir -p "$TEST_DIR"
echo -e "${BLUE}Created test directory: $TEST_DIR${RESET}"

# Copy example stackfile
cp "$EXAMPLES_DIR/simple_analysis.yaml" "$TEST_DIR/"
echo -e "${GREEN}Copied example stack to test directory${RESET}"

cd "$TEST_DIR"

echo -e "${BLUE}=== Creating Mock Agents ===${RESET}"

# Create mock agents
echo -e "${YELLOW}Creating text extractor agent...${RESET}"
cat > Sentinelfile-extractor << EOF
name: text-extractor
description: Extracts and processes text
version: 1.0.0
base_model: claude-3-sonnet-20240229
input:
  - name: text
    type: string
    description: "Text to extract from"
  - name: clean_text
    type: boolean
    default: true
    description: "Whether to clean the text"
output_format: "json"
system_prompt: |
  You are a specialized text extraction agent. Extract structured text.
prompt_template: |
  Extract text from the following input:
  
  {{text}}
  
  Clean text: {{clean_text}}
EOF

echo -e "${YELLOW}Building text extractor agent...${RESET}"
../sentinel build -t text-extractor:latest -f Sentinelfile-extractor .

echo -e "${YELLOW}Creating sentiment analyzer agent...${RESET}"
cat > Sentinelfile-sentiment << EOF
name: sentiment-analyzer
description: Analyzes sentiment in text
version: 1.0.0
base_model: claude-3-sonnet-20240229
input:
  - name: text
    type: string
    description: "Text to analyze"
  - name: model
    type: string
    default: "default"
    description: "Model to use for analysis"
output_format: "json"
system_prompt: |
  You are a specialized sentiment analysis agent. Determine sentiment.
prompt_template: |
  Analyze the sentiment of the following text:
  
  {{text}}
  
  Using model: {{model}}
EOF

echo -e "${YELLOW}Building sentiment analyzer agent...${RESET}"
../sentinel build -t sentiment-analyzer:latest -f Sentinelfile-sentiment .

echo -e "${YELLOW}Creating text summarizer agent...${RESET}"
cat > Sentinelfile-summarizer << EOF
name: text-summarizer
description: Summarizes text
version: 1.0.0
base_model: claude-3-sonnet-20240229
input:
  - name: text
    type: string
    description: "Text to summarize"
  - name: sentiment
    type: object
    description: "Sentiment analysis results"
  - name: style
    type: string
    default: "concise"
    description: "Summary style"
output_format: "json"
system_prompt: |
  You are a specialized text summarization agent. Create concise summaries.
prompt_template: |
  Summarize the following text:
  
  {{text}}
  
  Sentiment analysis: {{sentiment}}
  
  Summary style: {{style}}
EOF

echo -e "${YELLOW}Building text summarizer agent...${RESET}"
../sentinel build -t text-summarizer:latest -f Sentinelfile-summarizer .

echo -e "${GREEN}All mock agents built successfully${RESET}"

echo -e "${BLUE}=== Testing Stack List Command ===${RESET}"
../sentinel stack list

echo -e "${BLUE}=== Testing Stack Inspect Command ===${RESET}"
../sentinel stack inspect simple_analysis.yaml

echo -e "${BLUE}=== Testing Stack Run Command ===${RESET}"
echo "Creating test input..."
cat > input.txt << EOF
This is a sample text for testing. The SentinelStacks system is designed to make it easy to create and execute AI agent workflows. The stack engine orchestrates multiple agents to work together, passing data between them to accomplish complex tasks. This test demonstrates a simple analysis stack with text extraction, sentiment analysis, and summarization.
EOF

echo -e "${YELLOW}Running stack with test input...${RESET}"
../sentinel stack run -f simple_analysis.yaml --input="@input.txt" --verbose

echo -e "${BLUE}=== Testing Stack Init Command ===${RESET}"
../sentinel stack init custom-stack --template pipeline

echo -e "${YELLOW}Inspecting generated stackfile...${RESET}"
../sentinel stack inspect Stackfile.yaml

echo -e "${BLUE}=== Testing Natural Language Stack Creation ===${RESET}"
../sentinel stack init nlp-stack --nl "Create a pipeline with a text processor that extracts entities, a sentiment analyzer, and a report generator that combines the results"

echo -e "${YELLOW}Inspecting NL-generated stackfile...${RESET}"
../sentinel stack inspect Stackfile.yaml

echo -e "${GREEN}=== All tests completed successfully ===${RESET}"
echo -e "${BLUE}Test results directory: ${YELLOW}$TEST_DIR${RESET}"
