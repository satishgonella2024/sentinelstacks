package format

import (
	"fmt"
	"path/filepath"
	"strings"
)

// File extensions for SentinelStacks components
const (
	// AgentExtension is the extension for Sentinel Agent packages
	AgentExtension = ".agent.sntl"
	
	// StackExtension is the extension for Sentinel Stack packages
	StackExtension = ".stack.sntl"
	
	// AgentDefinitionExtension is the extension for Sentinel Agent definition files
	AgentDefinitionExtension = ".agent.yaml"
	
	// StackDefinitionExtension is the extension for Sentinel Stack definition files
	StackDefinitionExtension = ".stack.yaml"
	
	// SignatureExtension is the extension for detached signatures
	SignatureExtension = ".sig.sntl"
	
	// ManifestFileName is the name of the package manifest file
	ManifestFileName = "sentinel.manifest.json"
)

// FileFormatVersion is the current version of the file format specification
const FileFormatVersion = "1.0.0"

// MagicHeaders for identifying file types
var MagicHeaders = map[string][]byte{
	AgentExtension: []byte("SNTL-AGENT-PKG"),
	StackExtension: []byte("SNTL-STACK-PKG"),
}

// FormatInfo contains metadata about a file format
type FormatInfo struct {
	Extension       string
	Description     string
	VersionSupport  []string
	PrimaryMimeType string
}

// GetFormatInfo returns information about a file format based on its extension
func GetFormatInfo(filename string) *FormatInfo {
	ext := strings.ToLower(filepath.Ext(filename))
	fullExt := strings.ToLower(filepath.Ext(strings.TrimSuffix(filename, ext)) + ext)
	
	switch fullExt {
	case AgentExtension:
		return &FormatInfo{
			Extension:       AgentExtension,
			Description:     "Sentinel Agent Package",
			VersionSupport:  []string{"1.0.0"},
			PrimaryMimeType: "application/x-sentinel-agent",
		}
	case StackExtension:
		return &FormatInfo{
			Extension:       StackExtension,
			Description:     "Sentinel Stack Package",
			VersionSupport:  []string{"1.0.0"},
			PrimaryMimeType: "application/x-sentinel-stack",
		}
	case AgentDefinitionExtension:
		return &FormatInfo{
			Extension:       AgentDefinitionExtension,
			Description:     "Sentinel Agent Definition",
			VersionSupport:  []string{"1.0.0"},
			PrimaryMimeType: "application/x-sentinel-agent-def+yaml",
		}
	case StackDefinitionExtension:
		return &FormatInfo{
			Extension:       StackDefinitionExtension,
			Description:     "Sentinel Stack Definition",
			VersionSupport:  []string{"1.0.0"},
			PrimaryMimeType: "application/x-sentinel-stack-def+yaml",
		}
	case SignatureExtension:
		return &FormatInfo{
			Extension:       SignatureExtension,
			Description:     "Sentinel Signature File",
			VersionSupport:  []string{"1.0.0"},
			PrimaryMimeType: "application/x-sentinel-signature",
		}
	default:
		return nil
	}
}

// ValidateFormatVersion checks if a format version is supported
func ValidateFormatVersion(formatType, version string) error {
	switch formatType {
	case "agent", "stack", "signature":
		if version != "1.0.0" {
			return fmt.Errorf("unsupported %s format version: %s (supported: 1.0.0)", formatType, version)
		}
		return nil
	default:
		return fmt.Errorf("unknown format type: %s", formatType)
	}
}

// GetDefaultFilename returns a default filename for a given content
func GetDefaultFilename(name, version string, formatType string) string {
	switch formatType {
	case "agent":
		return fmt.Sprintf("%s-%s%s", name, version, AgentExtension)
	case "stack":
		return fmt.Sprintf("%s-%s%s", name, version, StackExtension)
	case "agent-def":
		return fmt.Sprintf("%s%s", name, AgentDefinitionExtension)
	case "stack-def":
		return fmt.Sprintf("%s%s", name, StackDefinitionExtension)
	case "signature":
		return fmt.Sprintf("%s-%s%s", name, version, SignatureExtension)
	default:
		return fmt.Sprintf("%s-%s.sntl", name, version)
	}
}

// GetFileType determines the file type from a filename
func GetFileType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	fullExt := strings.ToLower(filepath.Ext(strings.TrimSuffix(filename, ext)) + ext)
	
	switch fullExt {
	case AgentExtension:
		return "agent-package"
	case StackExtension:
		return "stack-package"
	case AgentDefinitionExtension:
		return "agent-definition"
	case StackDefinitionExtension:
		return "stack-definition"
	case SignatureExtension:
		return "signature"
	default:
		// Check for common filenames
		baseName := strings.ToLower(filepath.Base(filename))
		if baseName == "sentinelfile" || baseName == "sentinelfile.yaml" || baseName == "sentinelfile.yml" {
			return "agent-definition"
		}
		if baseName == "stackfile" || baseName == "stackfile.yaml" || baseName == "stackfile.yml" {
			return "stack-definition"
		}
		return "unknown"
	}
}
