#!/bin/bash

echo "Adding missing dependencies..."

# Add yaml.v2 dependency
go get gopkg.in/yaml.v2

# Add term dependency
go get golang.org/x/term

echo "Dependencies added!"
