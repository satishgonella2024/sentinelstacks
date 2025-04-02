#!/bin/bash
echo "Compiling SentinelStacks..."
go build -o sentinel main.go
chmod +x sentinel
echo "Done!"
