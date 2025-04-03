package packages

import (
	"archive/tar"
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

	"github.com/satishgonella2024/sentinelstacks/internal/registry/security"
	"github.com/satishgonella2024/sentinelstacks/pkg/types"
)

// FileType defines the type of file in a package
type FileType string

const (
	// FileTypeCode represents code files (e.g., .go, .py)
	FileTypeCode FileType = "code"
	
	// FileTypeConfig represents configuration files
	FileTypeConfig FileType = "config"
	
	// FileTypeManifest represents manifest files
	FileTypeManifest FileType = "manifest"
	
	// FileTypeDoc represents documentation files
	FileTypeDoc FileType = "doc"
	
	// FileTypeData represents data files
	FileTypeData FileType = "data"
)

// PackageManifest contains metadata about a package
type PackageManifest struct {
	Name          string                 `json:"name"`
	Version       string                 `json:"version"`
	Type          types.PackageType      `json:"type"`
	Description   string                 `json:"description"`
	Author        string                 `json:"author"`
	Created       time.Time              `json:"created"`
	Files         []FileInfo             `json:"files"`
	Dependencies  []types.Dependency     `json:"dependencies"`
	Signatures    []SignatureRecord      `json:"signatures"`
	SchemaVersion string                 `json:"schemaVersion"`
	Labels        map[string]string      `json:"labels,omitempty"`
}

// FileInfo represents a file in the package
type FileInfo struct {
	Path     string   `json:"path"`
	Size     int64    `json:"size"`
	SHA256   string   `json:"sha256"`
	IsMain   bool     `json:"isMain"`
	Type     FileType `json:"type"`
}

// SignatureRecord represents a stored signature
type SignatureRecord struct {
	KeyID     string    `json:"keyId"`
	Signer    string    `json:"signer"`
	Algorithm string    `json:"algorithm"`
	Signature string    `json:"signature"`
	Timestamp time.Time `json:"timestamp"`
}

// SentinelPackage is the standard format for distributing agents and stacks
type SentinelPackage struct {
	Manifest     PackageManifest
	Path         string
	keyManager   *security.KeyManager
	sourceFiles  map[string]string // Maps target paths to source paths
}

// NewSentinelPackage creates a new sentinel package
func NewSentinelPackage(pkgType types.PackageType, name, version, description, author string) *SentinelPackage {
	return &SentinelPackage{
		Manifest: PackageManifest{
			Name:          name,
			Version:       version,
			Type:          pkgType,
			Description:   description,
			Author:        author,
			Created:       time.Now().UTC(),
			Files:         []FileInfo{},
			Dependencies:  []types.Dependency{},
			Signatures:    []SignatureRecord{},
			SchemaVersion: "1.0",
			Labels:        make(map[string]string),
		},
		sourceFiles: make(map[string]string),
	}
}

// SetKeyManager sets the key manager for signing operations
func (p *SentinelPackage) SetKeyManager(km *security.KeyManager) {
	p.keyManager = km
}

// AddFile adds a file to the package
func (p *SentinelPackage) AddFile(sourcePath, targetPath string, isMain bool, fileType FileType) error {
	// Check that the source file exists
	fileInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	
	if fileInfo.IsDir() {
		return fmt.Errorf("cannot add directory as file: %s", sourcePath)
	}
	
	// Calculate SHA256 hash
	hash, err := calculateFileSHA256(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to calculate file hash: %w", err)
	}
	
	// Add to manifest
	p.Manifest.Files = append(p.Manifest.Files, FileInfo{
		Path:     targetPath,
		Size:     fileInfo.Size(),
		SHA256:   hash,
		IsMain:   isMain,
		Type:     fileType,
	})
	
	// Store source path for later packaging
	p.sourceFiles[targetPath] = sourcePath
	
	return nil
}

// AddDirectory adds all files in a directory to the package
func (p *SentinelPackage) AddDirectory(sourceDirPath, targetPrefix string) error {
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
		return p.AddFile(path, targetPath, false, fileType)
	})
}

// AddDependency adds a dependency to the package
func (p *SentinelPackage) AddDependency(name, version string, depType types.PackageType, required bool) {
	p.Manifest.Dependencies = append(p.Manifest.Dependencies, types.Dependency{
		Name:     name,
		Version:  version,
		Type:     depType,
		Required: required,
	})
}

// AddLabel adds a label to the package metadata
func (p *SentinelPackage) AddLabel(key, value string) {
	p.Manifest.Labels[key] = value
}

// Sign signs the package manifest
func (p *SentinelPackage) Sign(keyID, signer string) error {
	if p.keyManager == nil {
		return fmt.Errorf("key manager not set")
	}
	
	// Marshal the manifest without signatures
	tempManifest := p.Manifest
	tempManifest.Signatures = []SignatureRecord{}
	manifestBytes, err := json.Marshal(tempManifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	
	// Sign the manifest
	signature, err := p.keyManager.Sign(manifestBytes, keyID, signer)
	if err != nil {
		return fmt.Errorf("failed to sign manifest: %w", err)
	}
	
	// Add signature to manifest
	p.Manifest.Signatures = append(p.Manifest.Signatures, SignatureRecord{
		KeyID:     signature.Info.KeyID,
		Signer:    signature.Info.Signer,
		Algorithm: signature.Info.Algorithm,
		Signature: signature.Data,
		Timestamp: signature.Info.Timestamp,
	})
	
	return nil
}

// Package creates a package file containing all added files
func (p *SentinelPackage) Package(outputPath string) error {
	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()
	
	// Create gzip writer
	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()
	
	// Create tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()
	
	// Add manifest file
	manifestBytes, err := json.MarshalIndent(p.Manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}
	
	manifestHeader := &tar.Header{
		Name: "sentinel-manifest.json",
		Mode: 0644,
		Size: int64(len(manifestBytes)),
	}
	
	if err := tarWriter.WriteHeader(manifestHeader); err != nil {
		return fmt.Errorf("failed to write manifest header: %w", err)
	}
	
	if _, err := tarWriter.Write(manifestBytes); err != nil {
		return fmt.Errorf("failed to write manifest content: %w", err)
	}
	
	// Add all files from the manifest
	for _, fileInfo := range p.Manifest.Files {
		sourcePath, exists := p.sourceFiles[fileInfo.Path]
		if !exists {
			return fmt.Errorf("source path for %s not found", fileInfo.Path)
		}
		
		if err := addFileToTar(tarWriter, sourcePath, fileInfo.Path); err != nil {
			return fmt.Errorf("failed to add file to package: %w", err)
		}
	}
	
	// Store the package path
	p.Path = outputPath
	
	return nil
}

// Unpackage extracts a package to a directory
func (p *SentinelPackage) Unpackage(packagePath, outputDir string) error {
	// Open the package file
	file, err := os.Open(packagePath)
	if err != nil {
		return fmt.Errorf("failed to open package: %w", err)
	}
	defer file.Close()
	
	// Create gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()
	
	// Create tar reader
	tarReader := tar.NewReader(gzipReader)
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Extract all files
	var manifestFound bool
	
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar: %w", err)
		}
		
		// Handle different file types
		switch header.Typeflag {
		case tar.TypeReg, tar.TypeRegA:
			// Regular file
			outPath := filepath.Join(outputDir, header.Name)
			outDir := filepath.Dir(outPath)
			
			// Create directory if it doesn't exist
			if err := os.MkdirAll(outDir, 0755); err != nil {
				return fmt.Errorf("error creating directory: %w", err)
			}
			
			// Create file
			outFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("error creating file: %w", err)
			}
			
			// Copy content
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("error copying file content: %w", err)
			}
			outFile.Close()
			
			// Check if this is the manifest
			if header.Name == "sentinel-manifest.json" {
				manifestData, err := os.ReadFile(outPath)
				if err != nil {
					return fmt.Errorf("error reading manifest: %w", err)
				}
				
				if err := json.Unmarshal(manifestData, &p.Manifest); err != nil {
					return fmt.Errorf("error parsing manifest: %w", err)
				}
				
				manifestFound = true
			}
		}
	}
	
	if !manifestFound {
		return fmt.Errorf("no manifest found in package")
	}
	
	// Update package path
	p.Path = packagePath
	
	return nil
}

// VerifyIntegrity checks that all files match their recorded hashes
func (p *SentinelPackage) VerifyIntegrity(baseDir string) (bool, []string, error) {
	var failures []string
	
	for _, fileInfo := range p.Manifest.Files {
		filePath := filepath.Join(baseDir, fileInfo.Path)
		
		// Calculate file hash
		hash, err := calculateFileSHA256(filePath)
		if err != nil {
			return false, failures, fmt.Errorf("failed to calculate hash for %s: %w", fileInfo.Path, err)
		}
		
		// Compare with recorded hash
		if hash != fileInfo.SHA256 {
			failures = append(failures, fileInfo.Path)
		}
	}
	
	return len(failures) == 0, failures, nil
}

// VerifySignatures verifies all signatures on the package
func (p *SentinelPackage) VerifySignatures() (bool, []string, error) {
	if p.keyManager == nil {
		return false, nil, fmt.Errorf("key manager not set")
	}
	
	var validSigners []string
	
	// Create manifest copy without signatures
	tempManifest := p.Manifest
	tempManifest.Signatures = []SignatureRecord{}
	manifestBytes, err := json.Marshal(tempManifest)
	if err != nil {
		return false, nil, fmt.Errorf("failed to marshal manifest: %w", err)
	}
	
	// Verify each signature
	for _, sig := range p.Manifest.Signatures {
		// Create signature object
		signature := &security.Signature{
			Data: sig.Signature,
			Info: security.SignatureInfo{
				KeyID:     sig.KeyID,
				Signer:    sig.Signer,
				Algorithm: sig.Algorithm,
				Timestamp: sig.Timestamp,
			},
		}
		
		// Verify signature
		err := p.keyManager.Verify(signature, manifestBytes)
		if err == nil {
			validSigners = append(validSigners, sig.Signer)
		}
	}
	
	return len(validSigners) > 0, validSigners, nil
}

// GetMainFiles returns all main files in the package
func (p *SentinelPackage) GetMainFiles() []FileInfo {
	var mainFiles []FileInfo
	
	for _, fileInfo := range p.Manifest.Files {
		if fileInfo.IsMain {
			mainFiles = append(mainFiles, fileInfo)
		}
	}
	
	return mainFiles
}

// ToPackageInfo converts to a types.PackageInfo
func (p *SentinelPackage) ToPackageInfo() types.PackageInfo {
	tags := make([]string, 0, len(p.Manifest.Labels))
	for key, value := range p.Manifest.Labels {
		tags = append(tags, fmt.Sprintf("%s:%s", key, value))
	}

	// Create metadata map
	metadata := map[string]interface{}{
		"schemaVersion": p.Manifest.SchemaVersion,
		"fileCount":     len(p.Manifest.Files),
		"signed":        len(p.Manifest.Signatures) > 0,
	}

	return types.PackageInfo{
		Name:         p.Manifest.Name,
		Version:      p.Manifest.Version,
		Type:         p.Manifest.Type,
		Description:  p.Manifest.Description,
		Author:       p.Manifest.Author,
		CreatedAt:    p.Manifest.Created,
		UpdatedAt:    p.Manifest.Created, // Use created as updated for now
		Tags:         tags,
		Dependencies: p.Manifest.Dependencies,
		Metadata:     metadata,
	}
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

// determineFileType guesses the file type based on extension
func determineFileType(path string) FileType {
	ext := strings.ToLower(filepath.Ext(path))
	
	switch ext {
	case ".go", ".py", ".js", ".ts", ".rb", ".c", ".cpp", ".cs", ".java", ".php":
		return FileTypeCode
	case ".md", ".txt", ".rst", ".html", ".pdf":
		return FileTypeDoc
	case ".json", ".yaml", ".yml", ".toml", ".ini", ".config":
		return FileTypeConfig
	case ".csv", ".xlsx", ".xls", ".db", ".sqlite", ".bin":
		return FileTypeData
	default:
		// Check for specific filenames
		baseName := strings.ToLower(filepath.Base(path))
		if baseName == "sentinelfile" || baseName == "stackfile.yaml" || baseName == "stackfile.yml" {
			return FileTypeManifest
		}
		return FileTypeData
	}
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
