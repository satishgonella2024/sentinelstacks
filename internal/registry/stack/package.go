package stack

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/satishgonella2024/sentinelstacks/internal/stack"
)

// StackPackage represents a packaged stack for transport or storage
type StackPackage struct {
	Spec      stack.StackSpec  `json:"spec"`
	DependsOn []AgentReference `json:"dependsOn"`
	Files     []PackageFile    `json:"files"`
}

// AgentReference represents a reference to an agent used by the stack
type AgentReference struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

// PackageFile represents a file included in the stack package
type PackageFile struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	IsMain   bool   `json:"isMain"`
	Checksum string `json:"checksum"`
}

// PackageStack packages a stack for transport
func PackageStack(stackFilePath string, outputPath string) (*StackPackage, error) {
	// Read stack file
	content, err := ioutil.ReadFile(stackFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read stack file: %w", err)
	}

	// Parse stack file
	var spec stack.StackSpec
	switch {
	case strings.HasSuffix(stackFilePath, ".yaml") || strings.HasSuffix(stackFilePath, ".yml"):
		if err := yaml.Unmarshal(content, &spec); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	case strings.HasSuffix(stackFilePath, ".json"):
		if err := json.Unmarshal(content, &spec); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %s", filepath.Ext(stackFilePath))
	}

	// Extract agent dependencies
	var dependencies []AgentReference
	for _, agent := range spec.Agents {
		// Parse agent reference
		parts := strings.Split(agent.Uses, ":")
		name := parts[0]
		tag := "latest"
		if len(parts) > 1 {
			tag = parts[1]
		}

		// Add to dependencies
		dependencies = append(dependencies, AgentReference{
			Name: name,
			Tag:  tag,
		})
	}

	// Package files - start with the stack file
	stackDir := filepath.Dir(stackFilePath)
	files := []PackageFile{
		{
			Path:   filepath.Base(stackFilePath),
			Size:   int64(len(content)),
			IsMain: true,
		},
	}

	// Look for additional files
	additionalFiles, err := findAdditionalFiles(stackDir, spec)
	if err != nil {
		return nil, fmt.Errorf("failed to find additional files: %w", err)
	}
	files = append(files, additionalFiles...)

	// Create package
	pkg := &StackPackage{
		Spec:      spec,
		DependsOn: dependencies,
		Files:     files,
	}

	// Create archive if outputPath provided
	if outputPath != "" {
		if err := createPackageArchive(pkg, stackDir, outputPath); err != nil {
			return nil, fmt.Errorf("failed to create archive: %w", err)
		}
	}

	return pkg, nil
}

// UnpackageStack extracts a stack package
func UnpackageStack(packagePath string, outputDir string) (*stack.StackSpec, string, error) {
	// Open the archive file
	file, err := os.Open(packagePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open package: %w", err)
	}
	defer file.Close()

	// Create a gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Extract files
	var mainStackFile string
	var stackSpec stack.StackSpec
	var metadataFound bool

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, "", fmt.Errorf("failed to read tar header: %w", err)
		}

		// Handle different file types
		switch header.Typeflag {
		case tar.TypeReg: // Regular file
			// Create output file
			outPath := filepath.Join(outputDir, header.Name)
			outDir := filepath.Dir(outPath)
			if err := os.MkdirAll(outDir, 0755); err != nil {
				return nil, "", fmt.Errorf("failed to create directory: %w", err)
			}

			outFile, err := os.Create(outPath)
			if err != nil {
				return nil, "", fmt.Errorf("failed to create file: %w", err)
			}

			// Copy file contents
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return nil, "", fmt.Errorf("failed to write file: %w", err)
			}
			outFile.Close()

			// If this is the metadata file, parse it
			if header.Name == "metadata.json" {
				metadataBytes, err := ioutil.ReadFile(outPath)
				if err != nil {
					return nil, "", fmt.Errorf("failed to read metadata: %w", err)
				}

				var pkg StackPackage
				if err := json.Unmarshal(metadataBytes, &pkg); err != nil {
					return nil, "", fmt.Errorf("failed to parse metadata: %w", err)
				}

				stackSpec = pkg.Spec
				metadataFound = true

				// Find main stack file
				for _, file := range pkg.Files {
					if file.IsMain {
						mainStackFile = filepath.Join(outputDir, file.Path)
						break
					}
				}
			}
		}
	}

	// If metadata wasn't found, try to find a stackfile
	if !metadataFound {
		// Look for Stackfile.yaml
		candidates := []string{
			filepath.Join(outputDir, "Stackfile.yaml"),
			filepath.Join(outputDir, "Stackfile.yml"),
			filepath.Join(outputDir, "stack.yaml"),
			filepath.Join(outputDir, "stack.yml"),
		}

		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				mainStackFile = candidate
				// Parse stack file
				content, err := ioutil.ReadFile(candidate)
				if err != nil {
					return nil, "", fmt.Errorf("failed to read stack file: %w", err)
				}

				if err := yaml.Unmarshal(content, &stackSpec); err != nil {
					return nil, "", fmt.Errorf("failed to parse stack file: %w", err)
				}
				break
			}
		}
	}

	if mainStackFile == "" {
		return nil, "", fmt.Errorf("no stack file found in package")
	}

	return &stackSpec, mainStackFile, nil
}

// findAdditionalFiles looks for files referenced by the stack
func findAdditionalFiles(baseDir string, spec stack.StackSpec) ([]PackageFile, error) {
	var files []PackageFile

	// Check for a README
	readmeFiles := []string{"README.md", "README.txt", "README"}
	for _, readmeFile := range readmeFiles {
		readmePath := filepath.Join(baseDir, readmeFile)
		if _, err := os.Stat(readmePath); err == nil {
			fileInfo, err := os.Stat(readmePath)
			if err != nil {
				continue
			}
			files = append(files, PackageFile{
				Path:   readmeFile,
				Size:   fileInfo.Size(),
				IsMain: false,
			})
			break
		}
	}

	// Check for other common files (docs, examples, etc.)
	dirsToCheck := []string{"docs", "examples"}
	for _, dir := range dirsToCheck {
		dirPath := filepath.Join(baseDir, dir)
		if _, err := os.Stat(dirPath); err == nil {
			// Walk directory
			filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if !info.IsDir() {
					relPath, err := filepath.Rel(baseDir, path)
					if err != nil {
						return nil
					}
					files = append(files, PackageFile{
						Path:   relPath,
						Size:   info.Size(),
						IsMain: false,
					})
				}
				return nil
			})
		}
	}

	return files, nil
}

// createPackageArchive creates a compressed archive of the stack
func createPackageArchive(pkg *StackPackage, sourceDir, outputPath string) error {
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

	// Add metadata file
	metadataBytes, err := json.Marshal(pkg)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metadataHeader := &tar.Header{
		Name: "metadata.json",
		Mode: 0644,
		Size: int64(len(metadataBytes)),
	}
	if err := tarWriter.WriteHeader(metadataHeader); err != nil {
		return fmt.Errorf("failed to write metadata header: %w", err)
	}
	if _, err := tarWriter.Write(metadataBytes); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	// Add all files
	for _, file := range pkg.Files {
		filePath := filepath.Join(sourceDir, file.Path)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("failed to stat file: %w", err)
		}

		// Create tar header
		header, err := tar.FileInfoHeader(fileInfo, "")
		if err != nil {
			return fmt.Errorf("failed to create header: %w", err)
		}
		header.Name = file.Path

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}

		// Open and copy file content
		fileContent, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		if _, err := io.Copy(tarWriter, fileContent); err != nil {
			fileContent.Close()
			return fmt.Errorf("failed to write file content: %w", err)
		}
		fileContent.Close()
	}

	return nil
}
