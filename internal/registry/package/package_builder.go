package packages

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/registry/format"
	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/registry/security"
	"github.com/subrahmanyagonella/the-repo/sentinelstacks/internal/stack"
)

// PackageBuilder builds SentinelStacks packages in the standard format
type PackageBuilder struct {
	packageType  PackageType
	name         string
	version      string
	description  string
	author       string
	files        []FileEntry
	dependencies []Dependency
	labels       map[string]string
	keyManager   *security.KeyManager
	buildTime    time.Time
}

// FileEntry represents a file to include in the package
type FileEntry struct {
	SourcePath string
	TargetPath string
	IsMain     bool
	Type       FileType
}

// NewPackageBuilder creates a new package builder
func NewPackageBuilder(pkgType PackageType, name, version, description, author string) *PackageBuilder {
	return &PackageBuilder{
		packageType:  pkgType,
		name:         name,
		version:      version,
		description:  description,
		author:       author,
		files:        []FileEntry{},
		dependencies: []Dependency{},
		labels:       make(map[string]string),
		buildTime:    time.Now().UTC(),
	}
}

// SetKeyManager sets the key manager for signing operations
func (b *PackageBuilder) SetKeyManager(km *security.KeyManager) {
	b.keyManager = km
}

// AddFile adds a file to the package
func (b *PackageBuilder) AddFile(sourcePath, targetPath string, isMain bool, fileType FileType) error {
	// Check that the source file exists
	_, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("source file not found: %w", err)
	}
	
	// Add the file entry
	b.files = append(b.files, FileEntry{
		SourcePath: sourcePath,
		TargetPath: targetPath,
		IsMain:     isMain,
		Type:       fileType,
	})
	
	return nil
}

// AddLabel adds a label to the package
func (b *PackageBuilder) AddLabel(key, value string) {
	b.labels[key] = value
}

// AddDependency adds a dependency to the package
func (b *PackageBuilder) AddDependency(name, version string, depType PackageType, required bool) {
	b.dependencies = append(b.dependencies, Dependency{
		Name:     name,
		Version:  version,
		Type:     depType,
		Required: required,
	})
}

// AddDirectory recursively adds files from a directory
func (b *PackageBuilder) AddDirectory(sourceDirPath, targetPrefix string) error {
	return filepath.Walk(sourceDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories themselves
		if info.IsDir() {
			return nil
		}
		
		// Calculate relative path
		relPath, err := filepath.Rel(sourceDirPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		
		// Create target path
		targetPath := filepath.Join(targetPrefix, relPath)
		
		// Determine file type based on extension
		fileType := determineFileType(path)
		
		// Add file to package
		return b.AddFile(path, targetPath, false, fileType)
	})
}

// ImportStackDefinition imports a stack definition into the package
func (b *PackageBuilder) ImportStackDefinition(stackFilePath string) error {
	// Verify that this is a stack package
	if b.packageType != PackageTypeStack {
		return fmt.Errorf("cannot import stack definition into non-stack package")
	}
	
	// Read the stack file
	content, err := os.ReadFile(stackFilePath)
	if err != nil {
		return fmt.Errorf("failed to read stack file: %w", err)
	}
	
	// Parse the stack file
	var stackSpec stack.StackSpec
	if strings.HasSuffix(stackFilePath, ".json") {
		if err := json.Unmarshal(content, &stackSpec); err != nil {
			return fmt.Errorf("failed to parse stack JSON: %w", err)
		}
	} else {
		// Assume YAML by default
		// In a real implementation, this would use yaml.Unmarshal
		return fmt.Errorf("YAML parsing not implemented")
	}
	
	// Set package metadata from stack spec
	b.name = stackSpec.Name
	b.description = stackSpec.Description
	if stackSpec.Version != "" {
		b.version = stackSpec.Version
	}
	
	// Add stack file as main file
	targetPath := format.GetDefaultFilename(stackSpec.Name, "", "stack-def")
	err = b.AddFile(stackFilePath, targetPath, true, FileTypeManifest)
	if err != nil {
		return fmt.Errorf("failed to add stack file: %w", err)
	}
	
	// Extract agent dependencies
	for _, agent := range stackSpec.Agents {
		parts := strings.Split(agent.Uses, ":")
		agentName := parts[0]
		agentVersion := "latest"
		if len(parts) > 1 {
			agentVersion = parts[1]
		}
		
		b.AddDependency(agentName, agentVersion, PackageTypeAgent, true)
	}
	
	return nil
}

// Build creates the package file
func (b *PackageBuilder) Build(outputPath string) error {
	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()
	
	// Set up buffered writer
	bufWriter := bufio.NewWriter(outFile)
	defer bufWriter.Flush()
	
	// Write magic header based on package type
	var magicHeader []byte
	if b.packageType == PackageTypeAgent {
		magicHeader = format.MagicHeaders[format.AgentExtension]
	} else {
		magicHeader = format.MagicHeaders[format.StackExtension]
	}
	
	if _, err := bufWriter.Write(magicHeader); err != nil {
		return fmt.Errorf("failed to write magic header: %w", err)
	}
	
	// Write format version (4 bytes)
	if _, err := bufWriter.WriteString("1.0\x00"); err != nil {
		return fmt.Errorf("failed to write format version: %w", err)
	}
	
	// Create gzip writer
	gzipWriter, err := gzip.NewWriterLevel(bufWriter, gzip.BestCompression)
	if err != nil {
		return fmt.Errorf("failed to create gzip writer: %w", err)
	}
	defer gzipWriter.Close()
	
	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()
	
	// Prepare manifest
	manifest := PackageManifest{
		Name:          b.name,
		Version:       b.version,
		Type:          b.packageType,
		Description:   b.description,
		Author:        b.author,
		Created:       b.buildTime,
		Files:         []FileInfo{},
		Dependencies:  b.dependencies,
		Signatures:    []SignatureRecord{},
		SchemaVersion: format.FileFormatVersion,
		Labels:        b.labels,
	}
	
	// Add files to the tar archive and collect file info
	for _, fileEntry := range b.files {
		// Calculate file hash
		hash, err := calculateFileSHA256(fileEntry.SourcePath)
		if err != nil {
			return fmt.Errorf("failed to calculate hash for %s: %w", fileEntry.SourcePath, err)
		}
		
		// Get file size
		fileInfo, err := os.Stat(fileEntry.SourcePath)
		if err != nil {
			return fmt.Errorf("failed to stat file: %w", err)
		}
		
		// Add file to manifest
		manifest.Files = append(manifest.Files, FileInfo{
			Path:     fileEntry.TargetPath,
			Size:     fileInfo.Size(),
			SHA256:   hash,
			IsMain:   fileEntry.IsMain,
			Type:     fileEntry.Type,
		})
		
		// Add file to tar archive
		if err := addFileToTar(tarWriter, fileEntry.SourcePath, fileEntry.TargetPath); err != nil {
			return fmt.Errorf("failed to add file to archive: %w", err)
		}
	}
	
	// Sign the manifest if a key manager is provided
	if b.keyManager != nil {
		manifestBytes, err := json.Marshal(manifest)
		if err != nil {
			return fmt.Errorf("failed to marshal manifest: %w", err)
		}
		
		// Get default key ID (in a real implementation, this would be configurable)
		keyID := "default"
		signer := b.author
		
		// Sign the manifest
		signature, err := b.keyManager.Sign(manifestBytes, keyID, signer)
		if err != nil {
			return fmt.Errorf("failed to sign manifest: %w", err)
		}
		
		// Add signature to manifest
		manifest.Signatures = append(manifest.Signatures, SignatureRecord{
			KeyID:     signature.Info.KeyID,
			Signer:    signature.Info.Signer,
			Algorithm: signature.Info.Algorithm,
			Signature: signature.Data,
			Timestamp: signature.Info.Timestamp,
		})
	}
	
	// Add manifest file to tar archive
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	
	manifestHeader := &tar.Header{
		Name: format.ManifestFileName,
		Mode: 0644,
		Size: int64(len(manifestBytes)),
	}
	
	if err := tarWriter.WriteHeader(manifestHeader); err != nil {
		return fmt.Errorf("failed to write manifest header: %w", err)
	}
	
	if _, err := tarWriter.Write(manifestBytes); err != nil {
		return fmt.Errorf("failed to write manifest content: %w", err)
	}
	
	// Ensure all data is written
	if err := tarWriter.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}
	
	if err := gzipWriter.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}
	
	if err := bufWriter.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}
	
	return nil
}

// calculateFileSHA256 calculates SHA256 hash of a file
func calculateFileSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// addFileToTar adds a file to a tar archive
func addFileToTar(tarWriter *tar.Writer, sourcePath, targetPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	
	header := &tar.Header{
		Name:    targetPath,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}
	
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	
	if _, err := io.Copy(tarWriter, file); err != nil {
		return fmt.Errorf("failed to copy file data: %w", err)
	}
	
	return nil
}

// BuildFromStackSpec creates a package directly from a stack spec
func BuildFromStackSpec(spec stack.StackSpec, outputPath string, author string) error {
	// Create a new package builder
	builder := NewPackageBuilder(PackageTypeStack, spec.Name, spec.Version, spec.Description, author)
	
	// Create a temporary file for the stack definition
	tempDir, err := os.MkdirTemp("", "sentinel-stack-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Write stack definition to temp file
	stackDefPath := filepath.Join(tempDir, format.GetDefaultFilename(spec.Name, "", "stack-def"))
	stackDefBytes, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stack spec: %w", err)
	}
	
	if err := os.WriteFile(stackDefPath, stackDefBytes, 0644); err != nil {
		return fmt.Errorf("failed to write stack definition: %w", err)
	}
	
	// Import the stack definition
	if err := builder.ImportStackDefinition(stackDefPath); err != nil {
		return fmt.Errorf("failed to import stack definition: %w", err)
	}
	
	// Build the package
	if err := builder.Build(outputPath); err != nil {
		return fmt.Errorf("failed to build package: %w", err)
	}
	
	return nil
}
